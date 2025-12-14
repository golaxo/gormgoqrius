package sqlite

import (
	"fmt"
	"slices"
	"testing"

	"github.com/golaxo/goqrius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gormgoqrius/tests"

	"github.com/golaxo/gormgoqrius"
)

func TestSQLiteDefault(t *testing.T) {
	t.Parallel()

	for name, test := range tests.DefaultScenarios() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := getDB(t)

			err := db.AutoMigrate(&tests.User{})
			require.NoError(t, err)

			db.CreateInBatches(&test.ToMigrate, len(test.ToMigrate))
			require.NoError(t, err)

			e, err := goqrius.Parse(test.Input)
			require.NoError(t, err)

			whereClause := gormgoqrius.WhereClause(e)
			var actual []*tests.User
			tx := db
			if whereClause != nil {
				tx = tx.Where(whereClause)
			}
			tx = tx.Debug().Find(&actual)
			require.NoError(t, tx.Error)

			expected := make([]*tests.User, 0, len(test.ToMigrate))
			for _, user := range test.ToMigrate {
				if slices.Contains(test.ExpectedIDs, user.ID) {
					expected = append(expected, user)
				}
			}

			assert.ElementsMatch(t, actual, expected)
		})
	}
}

func getDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%q?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	return db
}

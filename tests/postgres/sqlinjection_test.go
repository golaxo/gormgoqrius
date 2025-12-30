package postgres

import (
	"testing"

	"github.com/golaxo/goqrius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/clause"
	"gormgoqrius/tests"

	"github.com/golaxo/gormgoqrius"
)

func TestSQLInjection(t *testing.T) {
	t.Parallel()

	tdt := map[string]struct {
		input     string
		toMigrate []*tests.User
	}{
		"delete inside value filter": {
			input:     "name eq 'DELETE * FROM users;'",
			toMigrate: tests.DefaultUsers(),
		},
		"drop table inside value": {
			input:     "surname eq 'drop table users'",
			toMigrate: tests.DefaultUsers(),
		},
	}
	for name, test := range tdt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, deferFunc := getDB(t)
			defer deferFunc()

			err := db.AutoMigrate(&tests.User{})
			require.NoError(t, err)

			db.CreateInBatches(&test.toMigrate, len(test.toMigrate))

			var whereClause clause.Expression
			e, err := goqrius.Parse(test.input)
			if err == nil {
				whereClause = gormgoqrius.WhereClause(e)
			}

			var actual []*tests.User
			tx := db
			if whereClause != nil {
				tx = tx.Where(whereClause)
			}
			tx = tx.Debug().Find(&actual)
			require.NoError(t, tx.Error)

			var count int64
			db.Model(&tests.User{}).Count(&count)

			assert.Equal(t, int64(len(test.toMigrate)), count)
		})
	}
}

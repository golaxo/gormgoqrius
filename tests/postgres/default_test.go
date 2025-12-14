package postgres

import (
	"slices"
	"testing"

	"github.com/golaxo/goqrius"
	"github.com/golaxo/gormgoqrius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gormgoqrius/tests"
)

func TestPostgresDefault(t *testing.T) {
	t.Parallel()

	for name, test := range tests.DefaultScenarios() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, deferFunc := getDB(t)
			defer deferFunc()

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

func getDB(t *testing.T) (db *gorm.DB, deferFunc func()) {
	t.Helper()

	ctx := t.Context()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	deferFunc = func() {
		testcontainers.CleanupContainer(t, postgresContainer)
		errTC := testcontainers.TerminateContainer(postgresContainer)
		require.NoError(t, errTC)
	}

	dsn, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("can't get connection string: %s", err)
	}

	db, err = gorm.Open(postgresDriver.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	return db, deferFunc
}

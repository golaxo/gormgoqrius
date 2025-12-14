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

func TestWithPostgres(t *testing.T) {
	t.Parallel()

	email := "test@test.com"
	tdt := map[string]struct {
		input       string
		toMigrate   []*tests.User
		expectedIDs []int
	}{
		"no filter": {
			toMigrate: []*tests.User{
				{ID: 1, Name: "John", Surname: "Doe", Age: 20, Email: &email},
				{ID: 2, Name: "Jane", Surname: "Doe", Age: 10},
				{ID: 3, Name: "Alice", Surname: "Smith", Age: 66},
				{ID: 4, Name: "Bob", Surname: "Smith", Age: 30},
			},
			expectedIDs: []int{1, 2, 3, 4},
		},
		"name eq 'John'": {
			input: "name eq 'John'",
			toMigrate: []*tests.User{
				{ID: 1, Name: "John", Surname: "Doe", Age: 20, Email: &email},
				{ID: 2, Name: "Jane", Surname: "Doe", Age: 10},
				{ID: 3, Name: "Alice", Surname: "Smith", Age: 66},
				{ID: 4, Name: "Bob", Surname: "Smith", Age: 30},
			},
			expectedIDs: []int{1},
		},
		"name ne 'John'": {
			input: "name ne 'John'",
			toMigrate: []*tests.User{
				{ID: 1, Name: "John", Surname: "Doe", Age: 20, Email: &email},
				{ID: 2, Name: "Jane", Surname: "Doe", Age: 10},
				{ID: 3, Name: "Alice", Surname: "Smith", Age: 66},
				{ID: 4, Name: "Bob", Surname: "Smith", Age: 30},
			},
			expectedIDs: []int{2, 3, 4},
		},
		"not surname eq 'Doe'": {
			input: "not surname eq 'Doe'",
			toMigrate: []*tests.User{
				{ID: 1, Name: "John", Surname: "Doe", Age: 20, Email: &email},
				{ID: 2, Name: "Jane", Surname: "Doe", Age: 10},
				{ID: 3, Name: "Alice", Surname: "Smith", Age: 66},
				{ID: 4, Name: "Bob", Surname: "Smith", Age: 30},
			},
			expectedIDs: []int{3, 4},
		},
		"email eq null": {
			input: "email eq null",
			toMigrate: []*tests.User{
				{ID: 1, Name: "John", Surname: "Doe", Age: 20, Email: &email},
				{ID: 2, Name: "Jane", Surname: "Doe", Age: 10},
				{ID: 3, Name: "Alice", Surname: "Smith", Age: 66},
				{ID: 4, Name: "Bob", Surname: "Smith", Age: 30},
			},
			expectedIDs: []int{2, 3, 4},
		},
		"email ne null": {
			input: "email ne null",
			toMigrate: []*tests.User{
				{ID: 1, Name: "John", Surname: "Doe", Age: 20, Email: &email},
				{ID: 2, Name: "Jane", Surname: "Doe", Age: 10},
				{ID: 3, Name: "Alice", Surname: "Smith", Age: 66},
				{ID: 4, Name: "Bob", Surname: "Smith", Age: 30},
			},
			expectedIDs: []int{1},
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
			require.NoError(t, err)

			e, err := goqrius.Parse(test.input)
			require.NoError(t, err)

			whereClause := gormgoqrius.WhereClause(e)
			var actual []*tests.User
			tx := db
			if whereClause != nil {
				tx = tx.Where(whereClause)
			}
			tx = tx.Debug().Find(&actual)
			require.NoError(t, tx.Error)

			expected := make([]*tests.User, 0, len(test.toMigrate))
			for _, user := range test.toMigrate {
				if slices.Contains(test.expectedIDs, user.ID) {
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

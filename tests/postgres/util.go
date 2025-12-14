package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

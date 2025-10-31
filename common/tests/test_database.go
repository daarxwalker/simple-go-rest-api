package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"gocourse/common/database"
	"gocourse/migrations"
)

func TruncateAllTables(c context.Context, db database.DB) error {
	_, execErr := db.Exec(c, `TRUNCATE TABLE users RESTART IDENTITY CASCADE`)
	if execErr != nil {
		return fmt.Errorf("failed to truncate all tables: %w", execErr)
	}
	return nil
}

func CreateTestDatabase(t testing.TB) (*pgxpool.Pool, func()) {
	port := "5432"
	dbName, dbUser, dbPassword := "test", "test", "test"
	container, err := testcontainers.GenericContainer(
		t.Context(), testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image: "postgres:18-alpine",
				Env: map[string]string{
					"POSTGRES_USER":     dbUser,
					"POSTGRES_PASSWORD": dbPassword,
					"POSTGRES_DB":       dbName,
				},
				ExposedPorts: []string{port + "/tcp"},
				Cmd:          []string{"postgres", "-c", "fsync=off"},
				WaitingFor:   wait.ForListeningPort(nat.Port(port + "/tcp")),
			},
			Started: true,
		},
	)
	require.NoError(t, err)
	containerPort, err := container.MappedPort(t.Context(), nat.Port(port))
	require.NoError(t, err)
	t.Cleanup(
		func() {
			_ = container.Terminate(t.Context())
		},
	)
	uri := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		containerPort.Port(),
		dbName,
	)
	source, createIOFSErr := iofs.New(migrations.FS, ".")
	if createIOFSErr != nil {
		t.Fatalf("create iofs migrations failed: %v", createIOFSErr)
	}
	db, err := database.Connect(uri)
	require.NoError(t, err)
	migrator, createMigratorErr := migrate.NewWithSourceInstance(
		"iofs",
		source,
		uri,
	)
	if createMigratorErr != nil || migrator == nil {
		t.Fatalf("create migrator failed: %v", createMigratorErr)
	}
	if migrateUpErr := migrator.Up(); migrateUpErr != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("migrate up failed: %v", migrateUpErr)
	}
	require.NoError(t, err)
	return db, db.Close
}

package testutils

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SQLiteContainer holds the container and connection details for an in-memory SQLite instance.
type SQLiteContainer struct {
	testcontainers.Container
	ConnectionString string // In-memory DSN for SQLite
}

// CreateSQLiteContainer sets up a containerized in-memory SQLite instance for testing.
func CreateSQLiteContainer(ctx context.Context) (*SQLiteContainer, error) {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "nouchka/sqlite3:latest",
			WaitingFor: wait.ForLog("SQLite version").WithOccurrence(1).
				WithStartupTimeout(5 * time.Second),
		},
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	return &SQLiteContainer{
		Container:        container,
		ConnectionString: ":memory:",
	}, nil
}

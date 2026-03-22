package data_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Broderick-Westrope/tetrigo/internal/data"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Errorf("closing test DB: %v", closeErr)
		}
	})

	err = data.EnsureTablesExist(context.Background(), db)
	require.NoError(t, err)

	return db
}

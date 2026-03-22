package data

import (
	"context"
	"database/sql"

	// Import the sqlite3 driver.
	_ "github.com/mattn/go-sqlite3"
)

func NewDB(ctx context.Context, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = EnsureTablesExist(ctx, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func EnsureTablesExist(ctx context.Context, db *sql.DB) error {
	// Leaderboard table
	_, err := db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS leaderboard 
(id INTEGER PRIMARY KEY, game_mode TEXT, name TEXT, time INTEGER, score INTEGER, lines INTEGER, level INTEGER)`,
	)
	if err != nil {
		return err
	}

	return nil
}

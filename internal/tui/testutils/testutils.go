package testutils

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stuttgart-things/sthings-tetris/internal/data"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/require"
)

func WaitForMsgOfType[T tea.Msg](t *testing.T, tm *teatest.TestModel, msgCh chan T, timeout time.Duration) {
	msg := tm.WaitForMsg(t, func(msg tea.Msg) bool {
		_, ok := msg.(T)
		return ok
	}, teatest.WithDuration(timeout))

	if msg == nil {
		t.Error("WaitForMsg returned nil")
		return
	}

	concreteMsg, ok := msg.(T)
	if !ok {
		var zeroValue T
		t.Errorf("Expected %T, got %T", zeroValue, msg)
		return
	}
	msgCh <- concreteMsg
}

func SetupInMemoryDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	err = data.EnsureTablesExist(db)
	require.NoError(t, err)
	return db
}

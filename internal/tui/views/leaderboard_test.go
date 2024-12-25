package views

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/Broderick-Westrope/tetrigo/internal/data"
	"github.com/Broderick-Westrope/tetrigo/internal/tui"
	"github.com/Broderick-Westrope/x/exp/teatest"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"
)

func TestLeaderboard_TableEntries(t *testing.T) {
	db := setupInMemoryDB(t)
	repo := data.NewLeaderboardRepository(db)

	tt := map[string]struct {
		count int
	}{
		"0 (empty)": {
			count: 0,
		},
		"3 (partial)": {
			count: 3,
		},
		"50 (overfull)": {
			count: 50,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			for i := range tc.count {
				_, err := repo.Save(&data.Score{
					GameMode: t.Name(),
					Name:     fmt.Sprintf("user-%d", i),
					Time:     time.Second * time.Duration(i*2),
					Score:    i * 100,
					Lines:    i,
					Level:    i + 2,
				})
				require.NoError(t, err)
			}

			m, err := NewLeaderboardModel(&tui.LeaderboardInput{
				GameMode: t.Name(),
			}, db)
			require.NoError(t, err)

			tm := teatest.NewTestModel(t, m)
			tm.Send(tea.Quit())

			outBytes := []byte(tm.FinalModel(t).View())
			teatest.RequireEqualOutput(t, outBytes)
		})
	}
}

func TestLeaderboard_NewEntryInEmptyTable(t *testing.T) {
	db := setupInMemoryDB(t)

	m, err := NewLeaderboardModel(&tui.LeaderboardInput{
		GameMode: t.Name(),
		NewEntry: &data.Score{
			GameMode: t.Name(),
			Name:     "user-new",
			Time:     time.Minute,
			Score:    1000,
			Lines:    2,
			Level:    3,
		},
	}, db)
	require.NoError(t, err)

	tm := teatest.NewTestModel(t, m)
	tm.Send(tea.Quit())

	outBytes := []byte(tm.FinalModel(t).View())
	teatest.RequireEqualOutput(t, outBytes)
}

func TestLeaderboard_KeyboardNavigation(t *testing.T) {
	db := setupInMemoryDB(t)
	repo := data.NewLeaderboardRepository(db)

	for i := range 50 {
		_, err := repo.Save(&data.Score{
			GameMode: t.Name(),
			Name:     fmt.Sprintf("user-%d", i),
			Time:     time.Second * time.Duration(i*2),
			Score:    i * 100,
			Lines:    i,
			Level:    i + 2,
		})
		require.NoError(t, err)
	}

	m, err := NewLeaderboardModel(&tui.LeaderboardInput{
		GameMode: t.Name(),
		NewEntry: &data.Score{
			GameMode: t.Name(),
			Name:     "user-new",
			Time:     time.Minute,
			Score:    2001,
			Lines:    2,
			Level:    3,
		},
	}, db)
	require.NoError(t, err)

	tm := teatest.NewTestModel(t, m)
	for range 3 {
		tm.Send(tea.KeyMsg{Type: tea.KeyDown})
		time.Sleep(10 * time.Millisecond)
	}
	tm.Send(tea.Quit())

	outBytes := []byte(tm.FinalModel(t).View())
	teatest.RequireEqualOutput(t, outBytes)
}

func setupInMemoryDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	err = data.EnsureTablesExist(db)
	require.NoError(t, err)
	return db
}

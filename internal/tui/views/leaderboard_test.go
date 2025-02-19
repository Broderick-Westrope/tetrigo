package views

import (
	"fmt"
	"testing"
	"time"

	"github.com/stuttgart-things/sthings-tetris/internal/data"
	"github.com/stuttgart-things/sthings-tetris/internal/tui"
	"github.com/stuttgart-things/sthings-tetris/internal/tui/testutils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/require"
)

func TestLeaderboard_TableEntries(t *testing.T) {
	db := testutils.SetupInMemoryDB(t)
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
	db := testutils.SetupInMemoryDB(t)

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
	db := testutils.SetupInMemoryDB(t)
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

func TestLeaderboard_SwitchModeMsg(t *testing.T) {
	db := testutils.SetupInMemoryDB(t)

	m, err := NewLeaderboardModel(&tui.LeaderboardInput{
		GameMode: t.Name(),
	}, db)
	require.NoError(t, err)
	tm := teatest.NewTestModel(t, m)

	switchModeMsgCh := make(chan tui.SwitchModeMsg, 1)
	go testutils.WaitForMsgOfType(t, tm, switchModeMsgCh, time.Second)

	// exit the leaderboard
	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
	time.Sleep(10 * time.Millisecond)

	// Wait for switch mode message with timeout
	select {
	case switchModeMsg := <-switchModeMsgCh:
		require.Equal(t, tui.ModeMenu, switchModeMsg.Target)

		_, ok := switchModeMsg.Input.(*tui.MenuInput)
		require.True(t, ok, "Expected %T, got %T", &tui.MenuInput{}, switchModeMsg.Input)

	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for switch mode message")
	}
}

package views

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stuttgart-things/sthings-tetris/internal/config"
	"github.com/stuttgart-things/sthings-tetris/internal/data"
	"github.com/stuttgart-things/sthings-tetris/internal/tui"
	"github.com/stuttgart-things/sthings-tetris/internal/tui/components"
	"github.com/stuttgart-things/sthings-tetris/internal/tui/testutils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSingle_InitialOutput(t *testing.T) {
	r := rand.New(rand.NewPCG(0, 0))

	m, err := NewSingleModel(
		&tui.SingleInput{
			Mode:     tui.ModeMarathon,
			Level:    1,
			Username: "testuser",
		},
		&config.Config{
			NextQueueLength: 0,
			GhostEnabled:    true,
			LockDownMode:    "",
			MaxLevel:        0,
			EndOnMaxLevel:   false,
			Theme:           config.DefaultTheme(),
			Keys:            config.DefaultKeys(),
		},
		WithRandSource(r),
	)
	require.NoError(t, err)
	tm := teatest.NewTestModel(t, m)

	tm.Send(tea.Quit())
	outBytes := []byte(tm.FinalModel(t).View())
	teatest.RequireEqualOutput(t, outBytes)
}

func TestSingle_Interaction(t *testing.T) {
	mockGameStopwatch := components.NewMockStopwatch(t)
	mockGameStopwatch.EXPECT().Init().Return(nil)
	mockGameStopwatch.EXPECT().Update(mock.Anything).Return(mockGameStopwatch, nil)
	mockGameStopwatch.EXPECT().Elapsed().Return(time.Duration(0))

	m, err := NewSingleModel(
		&tui.SingleInput{
			Mode:     tui.ModeMarathon,
			Level:    1,
			Username: "testuser",
		},
		&config.Config{
			NextQueueLength: 0,
			GhostEnabled:    true,
			LockDownMode:    "",
			MaxLevel:        0,
			EndOnMaxLevel:   false,
			Theme:           config.DefaultTheme(),
			Keys:            config.DefaultKeys(),
		},
		WithRandSource(rand.New(rand.NewPCG(0, 0))),
	)
	require.NoError(t, err)
	m.gameStopwatch = mockGameStopwatch
	tm := teatest.NewTestModel(t, m)

	// hold
	tm.Send(tea.KeyMsg{Type: tea.KeySpace})
	time.Sleep(10 * time.Millisecond)

	// left 3
	for range 3 {
		tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
		time.Sleep(10 * time.Millisecond)
	}

	// hard drop
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("w")})
	time.Sleep(10 * time.Millisecond)

	// right
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	time.Sleep(10 * time.Millisecond)

	// hard drop
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("w")})
	time.Sleep(10 * time.Millisecond)

	// hold
	tm.Send(tea.KeyMsg{Type: tea.KeySpace})
	time.Sleep(10 * time.Millisecond)

	// right 4
	for range 4 {
		tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
		time.Sleep(10 * time.Millisecond)
	}

	// hard drop
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("w")})
	time.Sleep(10 * time.Millisecond)

	tm.Send(tea.Quit())
	outBytes := []byte(tm.FinalModel(t, teatest.WithFinalTimeout(time.Second)).View())
	teatest.RequireEqualOutput(t, outBytes)
}

// TODO: The golden file of this test shows that the ghost is
// written over the top of the placed Tet at the skyline.
// TODO: Fix formatting for the output. The golden file shows
// that escape sequences of the background are not handled
// correctly when overlaying the modal.
func TestSingle_GameOverOutput(t *testing.T) {
	mockGameStopwatch := components.NewMockStopwatch(t)
	mockGameStopwatch.EXPECT().Init().Return(nil)
	mockGameStopwatch.EXPECT().Update(mock.Anything).Return(mockGameStopwatch, nil)
	mockGameStopwatch.EXPECT().Elapsed().Return(time.Duration(0))

	m, err := NewSingleModel(
		&tui.SingleInput{
			Mode:     tui.ModeMarathon,
			Level:    1,
			Username: "testuser",
		},
		&config.Config{
			NextQueueLength: 0,
			GhostEnabled:    true,
			LockDownMode:    "",
			MaxLevel:        0,
			EndOnMaxLevel:   false,
			Theme:           config.DefaultTheme(),
			Keys:            config.DefaultKeys(),
		},
		WithRandSource(rand.New(rand.NewPCG(0, 0))),
	)
	require.NoError(t, err)
	m.gameStopwatch = mockGameStopwatch
	tm := teatest.NewTestModel(t, m)

	// hard drop 12
	for range 12 {
		tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("w")})
		time.Sleep(10 * time.Millisecond)
	}

	tm.Send(tea.Quit())
	outBytes := []byte(tm.FinalModel(t, teatest.WithFinalTimeout(time.Second)).View())
	teatest.RequireEqualOutput(t, outBytes)
}

func TestSingle_PausedOutput(t *testing.T) {
	mockGameStopwatch := components.NewMockStopwatch(t)
	mockGameStopwatch.EXPECT().Init().Return(nil)
	mockGameStopwatch.EXPECT().Update(mock.Anything).Return(mockGameStopwatch, nil)
	mockGameStopwatch.EXPECT().Elapsed().Return(time.Duration(0))
	mockGameStopwatch.EXPECT().Toggle().Return(nil)

	m, err := NewSingleModel(
		&tui.SingleInput{
			Mode:     tui.ModeMarathon,
			Level:    1,
			Username: "testuser",
		},
		&config.Config{
			NextQueueLength: 0,
			GhostEnabled:    true,
			LockDownMode:    "",
			MaxLevel:        0,
			EndOnMaxLevel:   false,
			Theme:           config.DefaultTheme(),
			Keys:            config.DefaultKeys(),
		},
		WithRandSource(rand.New(rand.NewPCG(0, 0))),
	)
	require.NoError(t, err)
	m.gameStopwatch = mockGameStopwatch
	tm := teatest.NewTestModel(t, m)

	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
	time.Sleep(10 * time.Millisecond)

	tm.Send(tea.Quit())
	outBytes := []byte(tm.FinalModel(t, teatest.WithFinalTimeout(time.Second)).View())
	teatest.RequireEqualOutput(t, outBytes)
}

func TestSingle_GameOverSwitchModeMsg(t *testing.T) {
	m, err := NewSingleModel(
		&tui.SingleInput{
			Mode:     tui.ModeMarathon,
			Level:    1,
			Username: "testuser",
		},
		&config.Config{
			NextQueueLength: 0,
			GhostEnabled:    true,
			LockDownMode:    "",
			MaxLevel:        0,
			EndOnMaxLevel:   false,
			Theme:           config.DefaultTheme(),
			Keys:            config.DefaultKeys(),
		},
		WithRandSource(rand.New(rand.NewPCG(0, 0))),
	)
	require.NoError(t, err)
	tm := teatest.NewTestModel(t, m)

	switchModeMsgCh := make(chan tui.SwitchModeMsg, 1)
	go testutils.WaitForMsgOfType(t, tm, switchModeMsgCh, time.Second)

	// hard drop 12
	for range 12 {
		tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("w")})
		time.Sleep(10 * time.Millisecond)
	}

	// continue past game over message
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(10 * time.Millisecond)

	// Wait for switch mode message with timeout
	select {
	case switchModeMsg := <-switchModeMsgCh:
		require.Equal(t, tui.ModeLeaderboard, switchModeMsg.Target)

		leaderboardInput, ok := switchModeMsg.Input.(*tui.LeaderboardInput)
		require.True(t, ok, "Expected %T, got %T", &tui.LeaderboardInput{}, switchModeMsg.Input)

		assert.Equal(t, tui.ModeMarathon.String(), leaderboardInput.GameMode)
		assert.EqualValues(t, &data.Score{
			ID:       0,
			Rank:     0,
			GameMode: tui.ModeMarathon.String(),
			Name:     "testuser",
			Time:     time.Duration(0),
			Score:    230,
			Lines:    0,
			Level:    1,
		}, leaderboardInput.NewEntry)

	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for switch mode message")
	}
}

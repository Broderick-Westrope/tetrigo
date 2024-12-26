package views

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/Broderick-Westrope/tetrigo/internal/tui"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/components"
	"github.com/Broderick-Westrope/x/exp/teatest"
	tea "github.com/charmbracelet/bubbletea"
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

package views

import (
	"math/rand/v2"
	"testing"

	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/Broderick-Westrope/tetrigo/internal/tui"
	"github.com/Broderick-Westrope/x/exp/teatest"
	tea "github.com/charmbracelet/bubbletea"
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

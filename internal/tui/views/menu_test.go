package views

import (
	"strings"
	"testing"
	"time"

	"github.com/stuttgart-things/sthings-tetris/internal/tui"
	"github.com/stuttgart-things/sthings-tetris/internal/tui/testutils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMenu_Output(t *testing.T) {
	m := NewMenuModel(&tui.MenuInput{})
	tm := teatest.NewTestModel(t, m)

	// Input username
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("testuser")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(10 * time.Millisecond)

	// Select game mode
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(10 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(10 * time.Millisecond)

	// Select level
	for range 3 {
		tm.Send(tea.KeyMsg{Type: tea.KeyDown})
		time.Sleep(10 * time.Millisecond)
	}

	tm.Send(tea.Quit())
	outBytes := []byte(tm.FinalModel(t).View())
	teatest.RequireEqualOutput(t, outBytes)
}

func TestMenu_SwitchModeMsg(t *testing.T) {
	modeToMoveDownCount := func(mode tui.Mode) int {
		switch mode {
		case tui.ModeMarathon:
			return 0
		case tui.ModeSprint:
			return 1
		case tui.ModeUltra:
			return 2
		case tui.ModeMenu:
			fallthrough
		case tui.ModeLeaderboard:
			fallthrough
		default:
			return 0
		}
	}

	tt := map[string]struct {
		username string
		mode     tui.Mode
		level    int
	}{
		"marathon; level 1": {
			username: "testuser",
			mode:     tui.ModeMarathon,
			level:    1,
		},
		"sprint; level 3": {
			username: "Perry_Crona@hotmail.com",
			mode:     tui.ModeSprint,
			level:    3,
		},
		"ultra; level 15": {
			username: "testuser",
			mode:     tui.ModeUltra,
			level:    15,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			m := NewMenuModel(&tui.MenuInput{})
			tm := teatest.NewTestModel(t, m)

			switchModeMsgCh := make(chan tui.SwitchModeMsg, 1)
			go testutils.WaitForMsgOfType(t, tm, switchModeMsgCh, time.Second)

			// Wait for initial render
			var out string
			teatest.WaitForOutput(t, tm.Output(), func(bytes []byte) bool {
				out = string(bytes)
				return strings.Contains(out, "Username:")
			}, teatest.WithDuration(time.Second))

			// Verify expected content is present
			assert.Contains(t, out, "Username:")
			assert.Contains(t, out, "Game Mode:")
			assert.Contains(t, out, "Starting Level:")

			// Input username
			tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tc.username)})
			tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
			time.Sleep(10 * time.Millisecond)

			// Select game mode
			for range modeToMoveDownCount(tc.mode) {
				tm.Send(tea.KeyMsg{Type: tea.KeyDown})
				time.Sleep(10 * time.Millisecond)
			}
			tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
			time.Sleep(10 * time.Millisecond)

			// Select level
			for range tc.level - 1 {
				tm.Send(tea.KeyMsg{Type: tea.KeyDown})
				time.Sleep(10 * time.Millisecond)
			}
			tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
			time.Sleep(10 * time.Millisecond)

			// Wait for switch mode message with timeout
			select {
			case switchModeMsg := <-switchModeMsgCh:
				require.Equal(t, tc.mode, switchModeMsg.Target)

				singleInput, ok := switchModeMsg.Input.(*tui.SingleInput)
				require.True(t, ok, "Expected %T, got %T", &tui.SingleInput{}, switchModeMsg.Input)

				assert.Equal(t, tc.mode, singleInput.Mode)
				assert.Equal(t, tc.level, singleInput.Level)
				assert.Equal(t, tc.username, singleInput.Username)

			case <-time.After(time.Second):
				t.Fatal("Timeout waiting for switch mode message")
			}
		})
	}
}

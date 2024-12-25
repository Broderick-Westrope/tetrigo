package views

import (
	"testing"
	"time"

	"github.com/Broderick-Westrope/tetrigo/internal/tui"
	"github.com/Broderick-Westrope/x/exp/teatest"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMenuFormCompletion(t *testing.T) {
	m := NewMenuModel(&tui.MenuInput{})
	tm := teatest.NewTestModel(t, m)

	switchModeMsgCh := make(chan tui.SwitchModeMsg, 1)
	go waitForMsgOfType(t, tm, switchModeMsgCh, time.Second)

	// Wait for initial render
	time.Sleep(100 * time.Millisecond)
	var out string
	teatest.WaitForOutput(t, tm.Output(), func(bytes []byte) bool {
		out = string(bytes)
		return len(out) > 0
	}, teatest.WithDuration(time.Second))

	// Verify expected content is present
	assert.Contains(t, out, "Username:")
	assert.Contains(t, out, "Game Mode:")
	assert.Contains(t, out, "Starting Level:")

	// Input username
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("testuser")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Select Marathon mode
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Select level 1
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// TODO(fix): The following msg send should not be necessary. The SwitchModeMsg should be sent automatically when the form is completed.
	tm.Send(tea.KeyMsg{Type: tea.KeyUp})
	time.Sleep(100 * time.Millisecond)

	// Wait for switch mode message with timeout
	select {
	case switchModeMsg := <-switchModeMsgCh:
		require.Equal(t, tui.ModeMarathon, switchModeMsg.Target)
		singleInput, ok := switchModeMsg.Input.(*tui.SingleInput)
		require.True(t, ok, "Expected %T, got %T", &tui.SingleInput{}, switchModeMsg.Input)
		assert.Equal(t, tui.ModeMarathon, singleInput.Mode)
		assert.Equal(t, 1, singleInput.Level)
		assert.Equal(t, "testuser", singleInput.Username)

	default:
		t.Fatal("Timeout waiting for switch mode message")
	}
}

func waitForMsgOfType[T tea.Msg](t *testing.T, tm *teatest.TestModel, msgCh chan T, timeout time.Duration) {
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

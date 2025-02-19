package views

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/stuttgart-things/sthings-tetris/internal/data"
	"github.com/stuttgart-things/sthings-tetris/internal/tui"
)

var _ tea.Model = &LeaderboardModel{}

type LeaderboardModel struct {
	keys *leaderboardKeyMap
	help help.Model

	repo  *data.LeaderboardRepository
	table table.Model

	width  int
	height int
}

func NewLeaderboardModel(in *tui.LeaderboardInput, db *sql.DB) (*LeaderboardModel, error) {
	repo := data.NewLeaderboardRepository(db)

	var err error
	newEntryID := 0
	if in.NewEntry != nil {
		if in.NewEntry.Name == "" {
			in.NewEntry.Name = "Anonymous"
		}

		newEntryID, err = repo.Save(in.NewEntry)
		if err != nil {
			return nil, fmt.Errorf("saving new entry: %w", err)
		}
	}

	scores, err := repo.All(in.GameMode)
	if err != nil {
		return nil, fmt.Errorf("fetching scores: %w", err)
	}

	return &LeaderboardModel{
		keys:  defaultLeaderboardKeyMap(),
		help:  help.New(),
		repo:  repo,
		table: buildLeaderboardTable(scores, newEntryID),
	}, nil
}

func (m *LeaderboardModel) Init() tea.Cmd {
	return nil
}

func (m *LeaderboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Exit):
			return m, tui.SwitchModeCmd(tui.ModeMenu, tui.NewMenuInput())
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *LeaderboardModel) View() string {
	output := m.table.View() + "\n" + m.help.View(m.keys)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, output)
}

func buildLeaderboardTable(scores []data.Score, focusID int) table.Model {
	cols := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "Name", Width: 10},
		{Title: "Time", Width: 10},
		{Title: "Score", Width: 10},
		{Title: "Lines", Width: 5},
		{Title: "Level", Width: 5},
	}

	focusIndex := 0
	rows := make([]table.Row, len(scores))
	for i, s := range scores {
		if s.ID == focusID {
			focusIndex = i
		}

		rows[i] = table.Row{
			strconv.Itoa(s.Rank),
			s.Name,
			s.Time.String(),
			strconv.Itoa(s.Score),
			strconv.Itoa(s.Lines),
			strconv.Itoa(s.Level),
		}
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithStyles(s),
	)
	t.SetCursor(focusIndex)

	return t
}

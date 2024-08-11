package leaderboard

import (
	"database/sql"
	"strconv"

	"github.com/Broderick-Westrope/tetrigo/cmd/tui/common"
	"github.com/Broderick-Westrope/tetrigo/internal/data"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	keys *keyMap

	repo  *data.LeaderboardRepository
	table table.Model
}

func NewModel(in *common.LeaderboardInput, db *sql.DB, keys *common.Keys) (Model, error) {
	repo := data.NewLeaderboardRepository(db)

	if in.NewEntry != nil {
		if err := repo.Save(in.NewEntry); err != nil {
			return Model{}, err
		}
	}

	scores, err := repo.All(in.GameMode)
	if err != nil {
		return Model{}, err
	}

	m := Model{
		keys:  constructKeyMap(keys),
		repo:  repo,
		table: getLeaderboardTable(scores),
	}
	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Exit):
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.table.View()
}

func getLeaderboardTable(scores []data.Score) table.Model {
	cols := []table.Column{
		{Title: "Name", Width: 10},
		{Title: "Time", Width: 30},
		{Title: "Score", Width: 10},
		{Title: "Lines", Width: 5},
		{Title: "Level", Width: 5},
	}

	rows := make([]table.Row, len(scores))
	for i, s := range scores {
		rows[i] = table.Row{
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

	return table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithStyles(s),
	)
}

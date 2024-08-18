package leaderboard

import (
	"database/sql"
	"strconv"

	"github.com/Broderick-Westrope/tetrigo/cmd/tetrigo/common"
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

func NewModel(in *common.LeaderboardInput, db *sql.DB) (Model, error) {
	repo := data.NewLeaderboardRepository(db)

	var err error
	newEntryID := 0
	if in.NewEntry != nil {
		if in.NewEntry.Name == "" {
			in.NewEntry.Name = "Anonymous"
		}

		newEntryID, err = repo.Save(in.NewEntry)
		if err != nil {
			return Model{}, err
		}
	}

	scores, err := repo.All(in.GameMode)
	if err != nil {
		return Model{}, err
	}

	return Model{
		keys:  defaultKeyMap(),
		repo:  repo,
		table: getLeaderboardTable(scores, newEntryID),
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, m.keys.Exit) {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.table.View() + "\n"
}

func getLeaderboardTable(scores []data.Score, focusID int) table.Model {
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

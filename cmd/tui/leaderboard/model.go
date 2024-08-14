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

	var err error
	newEntryId := 0
	if in.NewEntry != nil {
		newEntryId, err = repo.Save(in.NewEntry)
		if err != nil {
			return Model{}, err
		}
	}

	scores, err := repo.All(in.GameMode)
	if err != nil {
		return Model{}, err
	}

	return Model{
		keys:  constructKeyMap(keys),
		repo:  repo,
		table: getLeaderboardTable(scores, newEntryId),
	}, nil
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
	return m.table.View() + "\n"
}

func processScores(scores []data.Score, focusIndex, maxCount, topCount int) []data.Score {
	if len(scores) < maxCount {
		return scores
	}

	separator := data.Score{Rank: -1}
	if focusIndex < maxCount {
		return append(scores[:maxCount-1], separator)
	}

	// Collect top scores and add a separator
	topScores := append(scores[:topCount], separator)

	// Calculate padding (number of scores to show surrounding the new entry)
	remainingCount := maxCount - len(topScores)
	quotient := remainingCount / 2
	remainder := remainingCount % 2

	upperPadding := quotient
	lowerPadding := quotient
	if remainder != 0 { // When topCount is even
		upperPadding += 1
	}

	// Collect the new score and X scores on either side (padding)
	upperBound := focusIndex + upperPadding
	lowerBound := focusIndex - lowerPadding
	if upperBound > len(scores) {
		upperBound = len(scores)
	}
	if lowerBound < topCount {
		lowerBound = topCount
	}
	surroundingScores := scores[lowerBound:upperBound]

	// Combine the two slices with a separator (use a special value for separator)
	return append(topScores, surroundingScores...)
}

func getLeaderboardTable(scores []data.Score, focusId int) table.Model {
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
		if s.ID == focusId {
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

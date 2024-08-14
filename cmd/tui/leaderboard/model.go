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
	scores = processScores(scores, 17, 11, 5)

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

func processScores(scores []data.Score, newEntryRank, maxCount, topCount int) []data.Score {
	newEntryIdx := newEntryRank - 1

	if newEntryRank <= maxCount {
		// Return top X scores (maxCount)
		if len(scores) < maxCount {
			return scores
		}
		return scores[:maxCount]
	}

	// Collect top scores and add a separator
	topScores := append(scores[:topCount], data.Score{Rank: -1})

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
	upperBound := newEntryIdx + upperPadding
	lowerBound := newEntryIdx - lowerPadding
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

func getLeaderboardTable(scores []data.Score) table.Model {
	cols := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "Name", Width: 10},
		{Title: "Time", Width: 20},
		{Title: "Score", Width: 20},
		{Title: "Lines", Width: 5},
		{Title: "Level", Width: 5},
	}

	rows := make([]table.Row, len(scores))
	for i, s := range scores {
		if s.Rank == -1 {
			// Add a separator row
			rows = append(rows, table.Row{"...", "...", "...", "...", "...", "..."})
			continue
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

	return table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithStyles(s),
	)
}

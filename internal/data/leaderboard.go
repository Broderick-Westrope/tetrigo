package data

import (
	"database/sql"
	"time"
)

type Score struct {
	ID       int
	GameMode string
	Name     string
	Time     time.Duration
	Score    int
	Lines    int
	Level    int
}

type LeaderboardRepository struct {
	db *sql.DB
}

func NewLeaderboardRepository(db *sql.DB) *LeaderboardRepository {
	return &LeaderboardRepository{db}
}

func (r *LeaderboardRepository) All(gameMode string) ([]Score, error) {
	rows, err := r.db.Query("SELECT * FROM leaderboard WHERE game_mode = $1 ORDER BY score DESC, time ASC", gameMode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []Score
	for rows.Next() {
		var h Score
		if err := rows.Scan(&h.ID, &h.GameMode, &h.Name, &h.Time, &h.Score, &h.Lines, &h.Level); err != nil {
			return nil, err
		}
		scores = append(scores, h)
	}

	return scores, nil
}

func (r *LeaderboardRepository) Save(score *Score) error {
	_, err := r.db.Exec("INSERT INTO leaderboard (game_mode, name, time, score, lines, level) VALUES ($1, $2, $3, $4, $5, $6)", score.GameMode, score.Name, score.Time, score.Score, score.Lines, score.Level)
	return err
}

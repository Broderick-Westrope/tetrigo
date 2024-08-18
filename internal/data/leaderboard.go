package data

import (
	"database/sql"
	"time"
)

type Score struct {
	ID       int
	Rank     int
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
	if err = rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []Score
	for rows.Next() {
		var s Score
		if err = rows.Scan(&s.ID, &s.GameMode, &s.Name, &s.Time, &s.Score, &s.Lines, &s.Level); err != nil {
			return nil, err
		}
		s.Rank = len(scores) + 1
		scores = append(scores, s)
	}

	return scores, nil
}

// Save saves a score to the leaderboard and returns the ID of the new score.
func (r *LeaderboardRepository) Save(score *Score) (int, error) {
	res, err := r.db.Exec(
		"INSERT INTO leaderboard (game_mode, name, time, score, lines, level) VALUES ($1, $2, $3, $4, $5, $6)",
		score.GameMode, score.Name, score.Time, score.Score, score.Lines, score.Level,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), err
}

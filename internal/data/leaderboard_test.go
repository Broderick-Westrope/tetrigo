package data_test

import (
	"testing"
	"time"

	"github.com/Broderick-Westrope/tetrigo/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	t.Parallel()

	t.Run("returns non-zero ID", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := data.NewLeaderboardRepository(db)

		id, err := repo.Save(&data.Score{
			GameMode: "marathon",
			Name:     "Alice",
			Time:     5 * time.Minute,
			Score:    1000,
			Lines:    10,
			Level:    3,
		})
		require.NoError(t, err)
		assert.NotZero(t, id)
	})

	t.Run("multiple saves return incrementing IDs", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := data.NewLeaderboardRepository(db)

		id1, err := repo.Save(&data.Score{
			GameMode: "marathon",
			Name:     "Alice",
			Time:     5 * time.Minute,
			Score:    1000,
			Lines:    10,
			Level:    3,
		})
		require.NoError(t, err)

		id2, err := repo.Save(&data.Score{
			GameMode: "marathon",
			Name:     "Bob",
			Time:     3 * time.Minute,
			Score:    2000,
			Lines:    20,
			Level:    5,
		})
		require.NoError(t, err)

		assert.Greater(t, id2, id1, "second ID should be greater than first")
	})

	t.Run("full field round-trip including time.Duration", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := data.NewLeaderboardRepository(db)

		input := &data.Score{
			GameMode: "sprint",
			Name:     "Charlie",
			Time:     2*time.Minute + 30*time.Second + 500*time.Millisecond,
			Score:    5000,
			Lines:    40,
			Level:    7,
		}

		id, err := repo.Save(input)
		require.NoError(t, err)

		scores, err := repo.All("sprint")
		require.NoError(t, err)
		require.Len(t, scores, 1)

		got := scores[0]
		assert.Equal(t, id, got.ID)
		assert.Equal(t, input.GameMode, got.GameMode)
		assert.Equal(t, input.Name, got.Name)
		assert.Equal(t, input.Time, got.Time, "time.Duration should survive round-trip as nanoseconds")
		assert.Equal(t, input.Score, got.Score)
		assert.Equal(t, input.Lines, got.Lines)
		assert.Equal(t, input.Level, got.Level)
		assert.Equal(t, 1, got.Rank)
	})

	t.Run("save with zero score", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := data.NewLeaderboardRepository(db)

		id, err := repo.Save(&data.Score{
			GameMode: "marathon",
			Name:     "Newbie",
			Time:     10 * time.Second,
			Score:    0,
			Lines:    0,
			Level:    1,
		})
		require.NoError(t, err)
		assert.NotZero(t, id)

		scores, err := repo.All("marathon")
		require.NoError(t, err)
		require.Len(t, scores, 1)
		assert.Equal(t, 0, scores[0].Score)
	})
}

func TestAll(t *testing.T) {
	t.Parallel()

	t.Run("empty table returns empty slice", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := data.NewLeaderboardRepository(db)

		scores, err := repo.All("marathon")
		require.NoError(t, err)
		assert.Empty(t, scores)
	})

	t.Run("single entry has correct fields and rank 1", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := data.NewLeaderboardRepository(db)

		_, err := repo.Save(&data.Score{
			GameMode: "marathon",
			Name:     "Alice",
			Time:     5 * time.Minute,
			Score:    1000,
			Lines:    10,
			Level:    3,
		})
		require.NoError(t, err)

		scores, err := repo.All("marathon")
		require.NoError(t, err)
		require.Len(t, scores, 1)

		got := scores[0]
		assert.Equal(t, 1, got.Rank)
		assert.Equal(t, "marathon", got.GameMode)
		assert.Equal(t, "Alice", got.Name)
		assert.Equal(t, 5*time.Minute, got.Time)
		assert.Equal(t, 1000, got.Score)
		assert.Equal(t, 10, got.Lines)
		assert.Equal(t, 3, got.Level)
	})

	t.Run("multiple entries ordered by score DESC then time ASC", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := data.NewLeaderboardRepository(db)

		entries := []data.Score{
			{GameMode: "marathon", Name: "Low", Time: 1 * time.Minute, Score: 500, Lines: 5, Level: 2},
			{GameMode: "marathon", Name: "High", Time: 3 * time.Minute, Score: 3000, Lines: 30, Level: 8},
			{GameMode: "marathon", Name: "Mid", Time: 2 * time.Minute, Score: 1500, Lines: 15, Level: 5},
		}
		for i := range entries {
			_, err := repo.Save(&entries[i])
			require.NoError(t, err)
		}

		scores, err := repo.All("marathon")
		require.NoError(t, err)
		require.Len(t, scores, 3)

		// Ordered by score DESC: High (3000), Mid (1500), Low (500)
		assert.Equal(t, "High", scores[0].Name)
		assert.Equal(t, "Mid", scores[1].Name)
		assert.Equal(t, "Low", scores[2].Name)
	})

	t.Run("ranks are sequential with no gaps", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := data.NewLeaderboardRepository(db)

		for i := 1; i <= 5; i++ {
			_, err := repo.Save(&data.Score{
				GameMode: "marathon",
				Name:     "Player",
				Time:     time.Duration(i) * time.Minute,
				Score:    i * 100,
				Lines:    i,
				Level:    i,
			})
			require.NoError(t, err)
		}

		scores, err := repo.All("marathon")
		require.NoError(t, err)
		require.Len(t, scores, 5)

		for i, s := range scores {
			assert.Equal(t, i+1, s.Rank, "rank at index %d should be %d", i, i+1)
		}
	})

	t.Run("identical scores get different sequential ranks", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := data.NewLeaderboardRepository(db)

		// Two entries with same score but different times
		_, err := repo.Save(&data.Score{
			GameMode: "marathon", Name: "Tied1", Time: 2 * time.Minute, Score: 1000, Lines: 10, Level: 3,
		})
		require.NoError(t, err)
		_, err = repo.Save(&data.Score{
			GameMode: "marathon", Name: "Tied2", Time: 3 * time.Minute, Score: 1000, Lines: 10, Level: 3,
		})
		require.NoError(t, err)

		scores, err := repo.All("marathon")
		require.NoError(t, err)
		require.Len(t, scores, 2)

		assert.Equal(t, 1, scores[0].Rank)
		assert.Equal(t, 2, scores[1].Rank)
	})

	t.Run("game mode filtering: scores from other modes excluded", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := data.NewLeaderboardRepository(db)

		_, err := repo.Save(&data.Score{
			GameMode: "modeA", Name: "PlayerA", Time: time.Minute, Score: 1000, Lines: 10, Level: 3,
		})
		require.NoError(t, err)
		_, err = repo.Save(&data.Score{
			GameMode: "modeB", Name: "PlayerB", Time: time.Minute, Score: 2000, Lines: 20, Level: 5,
		})
		require.NoError(t, err)

		scoresA, err := repo.All("modeA")
		require.NoError(t, err)
		require.Len(t, scoresA, 1)
		assert.Equal(t, "PlayerA", scoresA[0].Name)

		scoresB, err := repo.All("modeB")
		require.NoError(t, err)
		require.Len(t, scoresB, 1)
		assert.Equal(t, "PlayerB", scoresB[0].Name)
	})

	t.Run("tie-breaking: same score, shorter time ranks higher", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := data.NewLeaderboardRepository(db)

		_, err := repo.Save(&data.Score{
			GameMode: "marathon", Name: "Slow", Time: 5 * time.Minute, Score: 1000, Lines: 10, Level: 3,
		})
		require.NoError(t, err)
		_, err = repo.Save(&data.Score{
			GameMode: "marathon", Name: "Fast", Time: 2 * time.Minute, Score: 1000, Lines: 10, Level: 3,
		})
		require.NoError(t, err)

		scores, err := repo.All("marathon")
		require.NoError(t, err)
		require.Len(t, scores, 2)

		// Same score, ordered by time ASC: Fast (2min) before Slow (5min)
		assert.Equal(t, "Fast", scores[0].Name)
		assert.Equal(t, "Slow", scores[1].Name)
	})

	t.Run("nonexistent game mode returns empty result", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := data.NewLeaderboardRepository(db)

		// Save something to another mode to ensure DB isn't empty
		_, err := repo.Save(&data.Score{
			GameMode: "marathon", Name: "Alice", Time: time.Minute, Score: 1000, Lines: 10, Level: 3,
		})
		require.NoError(t, err)

		scores, err := repo.All("nonexistent")
		require.NoError(t, err)
		assert.Empty(t, scores)
	})
}

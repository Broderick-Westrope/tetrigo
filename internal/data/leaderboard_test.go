package data

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Broderick-Westrope/tetrigo/internal/data/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type LeaderboardTestSuite struct {
	suite.Suite
	container *testutils.SQLiteContainer
	repo      *LeaderboardRepository
	ctx       context.Context
}

func (suite *LeaderboardTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	container, err := testutils.CreateSQLiteContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.container = container
	db, err := NewDB(container.ConnectionString) // Uses :memory:
	if err != nil {
		log.Fatal(err)
	}
	suite.repo = NewLeaderboardRepository(db)
}

func (suite *LeaderboardTestSuite) TearDownSuite() {
	err := suite.repo.db.Close()
	if err != nil {
		log.Fatalf("error closing db: %s", err)
	}
	err = suite.container.Terminate(suite.ctx)
	if err != nil {
		log.Fatalf("error terminating SQLite container: %s", err)
	}
}

func (suite *LeaderboardTestSuite) TestSaveScore() {
	t := suite.T()

	score := &Score{
		GameMode: "marathon",
		Name:     "TestPlayer",
		Time:     100 * time.Second,
		Score:    1000,
		Lines:    10,
		Level:    5,
	}
	savedID, err := suite.repo.Save(score)
	require.NoError(t, err)
	assert.Greater(t, savedID, 1)
}

func (suite *LeaderboardTestSuite) TestAllScores() {
	t := suite.T()

	scores := []*Score{
		{
			GameMode: "marathon",
			Name:     "player1",
			Time:     10 * time.Second,
			Score:    100,
			Lines:    1,
			Level:    1,
		},
		{
			GameMode: "marathon",
			Name:     "player2",
			Time:     20 * time.Second,
			Score:    200,
			Lines:    2,
			Level:    2,
		},
		{
			GameMode: "sprint",
			Name:     "player3",
			Time:     30 * time.Second,
			Score:    300,
			Lines:    3,
			Level:    3,
		},
	}

	for _, score := range scores {
		_, err := suite.repo.Save(score)
		require.NoError(t, err)
	}

	// Retrieve all scores
	results, err := suite.repo.All("marathon")
	require.NoError(t, err)
	require.Len(t, results, 2)

	testScoreEqual := func(s1, s2 Score) {
		assert.Equal(t, s1.GameMode, s2.GameMode)
		assert.Equal(t, s1.Name, s2.Name)
		assert.Equal(t, s1.Time, s2.Time)
		assert.Equal(t, s1.Score, s2.Score)
		assert.Equal(t, s1.Lines, s2.Lines)
		assert.Equal(t, s1.Level, s2.Level)
	}

	// order by score DESC, time ASC
	testScoreEqual(*scores[0], results[1])
	testScoreEqual(*scores[1], results[0])
}

func TestLeaderboardTestSuite(t *testing.T) {
	suite.Run(t, new(LeaderboardTestSuite))
}

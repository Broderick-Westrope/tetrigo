package leaderboard

import (
	"testing"

	"github.com/Broderick-Westrope/tetrigo/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestGetProcessedScores(t *testing.T) {
	tt := map[string]struct {
		scores       []data.Score
		newEntryRank int
		maxCount     int
		topCount     int
		want         []data.Score
	}{
		"less scores than maxCount": {
			scores: []data.Score{
				{Rank: 1},
				{Rank: 2},
				{Rank: 3},
			},
			newEntryRank: 1,
			maxCount:     5,
			want: []data.Score{
				{Rank: 1},
				{Rank: 2},
				{Rank: 3},
			},
		},
		"new entry is within topCount": {
			scores: []data.Score{
				{Rank: 1},
				{Rank: 2},
				{Rank: 3},
				{Rank: 4},
				{Rank: 5},
			},
			newEntryRank: 1,
			maxCount:     3,
			topCount:     3,
			want: []data.Score{
				{Rank: 1},
				{Rank: 2},
				{Rank: 3},
			},
		},
		"new entry is outside topCount; uneven topCount": {
			scores: []data.Score{
				{Rank: 1},
				{Rank: 2},
				{Rank: 3},
				{Rank: 4},
				{Rank: 5},
				{Rank: 6},
				{Rank: 7},
				{Rank: 8},
				{Rank: 9},
				{Rank: 10},
				{Rank: 11},
				{Rank: 12},
			},
			newEntryRank: 9,
			maxCount:     8,
			topCount:     3,
			want: []data.Score{
				{Rank: 1},
				{Rank: 2},
				{Rank: 3},
				{Rank: -1},
				{Rank: 7},
				{Rank: 8},
				{Rank: 9},
				{Rank: 10},
			},
		},
		"new entry is outside topCount; even topCount": {
			scores: []data.Score{
				{Rank: 1},
				{Rank: 2},
				{Rank: 3},
				{Rank: 4},
				{Rank: 5},
				{Rank: 6},
				{Rank: 7},
				{Rank: 8},
				{Rank: 9},
				{Rank: 10},
				{Rank: 11},
				{Rank: 12},
			},
			newEntryRank: 9,
			maxCount:     8,
			topCount:     2,
			want: []data.Score{
				{Rank: 1},
				{Rank: 2},
				{Rank: -1},
				{Rank: 7},
				{Rank: 8},
				{Rank: 9},
				{Rank: 10},
				{Rank: 11},
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result := processScores(tc.scores, tc.newEntryRank, tc.maxCount, tc.topCount)

			assert.EqualValues(t, tc.want, result)
		})
	}
}

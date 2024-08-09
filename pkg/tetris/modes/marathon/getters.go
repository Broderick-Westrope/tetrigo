package marathon

import (
	"time"

	"github.com/Broderick-Westrope/tetrigo/pkg/tetris"
)

func (g *Game) IsGameOver() bool {
	return g.gameOver
}

func (g *Game) GetVisibleMatrix() [][]byte {
	return g.matrix.GetVisible()
}

func (g *Game) GetBagTetriminos() []tetris.Tetrimino {
	return g.bag.GetElements()
}

func (g *Game) GetHoldTetrimino() *tetris.Tetrimino {
	return g.holdTet
}

func (g *Game) GetTotalScore() uint {
	return g.scoring.Total()
}

func (g *Game) GetLevel() uint {
	return g.scoring.Level()
}

func (g *Game) GetLinesCleared() uint {
	return g.scoring.Lines()
}

func (g *Game) GetDefaultFallInterval() time.Duration {
	return g.fall.DefaultInterval
}

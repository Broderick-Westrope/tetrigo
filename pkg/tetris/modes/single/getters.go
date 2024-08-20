package single

import (
	"time"

	"github.com/Broderick-Westrope/tetrigo/pkg/tetris"
)

func (g *Game) IsGameOver() bool {
	return g.gameOver
}

func (g *Game) GetVisibleMatrix() (tetris.Matrix, error) {
	matrix := g.matrix.DeepCopy()

	if g.ghostTet != nil {
		err := matrix.AddTetrimino(g.ghostTet)
		if err != nil {
			return nil, err
		}
	}

	if err := matrix.AddTetrimino(g.tetInPlay); err != nil {
		return nil, err
	}

	return matrix.GetVisible(), nil
}

func (g *Game) GetBagTetriminos() []tetris.Tetrimino {
	return g.nextQueue.GetElements()
}

func (g *Game) GetHoldTetrimino() *tetris.Tetrimino {
	return g.holdQueue
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

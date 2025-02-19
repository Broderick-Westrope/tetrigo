package single

import (
	"strconv"
	"sync/atomic"
	"time"

	"github.com/stuttgart-things/sthings-tetris/pkg/tetris"
)

var cacheNumber int64

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

func (g *Game) GetTotalScore() int {
	return g.scoring.Total()
}

func (g *Game) GetLevel() int {
	return g.scoring.Level()
}

func (g *Game) GetLinesCleared() int {
	return g.scoring.Lines()
}

func (g *Game) GetMessage(score int) string {

	var returnValue string

	if score == 0 {
		atomic.StoreInt64(&cacheNumber, int64(score))
		returnValue = "0"
	}

	if score != 0 {
		oldValue := atomic.LoadInt64(&cacheNumber)

		lines := int(oldValue) - score
		atomic.StoreInt64(&cacheNumber, int64(score))

		returnValue = strconv.Itoa(lines)
		atomic.StoreInt64(&cacheNumber, int64(score))

	}

	return returnValue

}

func (g *Game) GetDefaultFallInterval() time.Duration {
	return g.fall.DefaultInterval
}

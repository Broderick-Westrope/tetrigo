package marathon

import (
	"fmt"
	"time"

	"github.com/Broderick-Westrope/tetrigo/pkg/tetris"
)

type Game struct {
	matrix           tetris.Matrix     // The Matrix of cells on which the game is played
	nextQueue        *tetris.NextQueue // The queue of upcoming Tetriminos
	tetInPlay        *tetris.Tetrimino // The current Tetrimino in play
	holdQueue        *tetris.Tetrimino // The Tetrimino that is being held
	canHold          bool              // Whether the player can hold the current Tetrimino
	gameOver         bool              // Whether the game is over
	softDropStartRow int               // Records where the user began soft drop
	scoring          *tetris.Scoring   // The scoring system
	fall             *tetris.Fall      // The system for calculating the fall speed
}

func NewGame(level, maxLevel uint) (*Game, error) {
	matrix, err := tetris.NewMatrix(40, 10)
	if err != nil {
		return nil, err
	}
	nq := tetris.NewNextQueue(matrix.GetSkyline())

	game := &Game{
		matrix:    matrix,
		nextQueue: nq,
		tetInPlay: nq.Next(),
		holdQueue: tetris.EmptyTetrimino,
		canHold:   true,
		gameOver:  false,
		// TODO: is start line needed?
		softDropStartRow: matrix.GetHeight(),
		scoring:          tetris.NewScoring(level, maxLevel),
		fall:             tetris.NewFall(level),
	}

	// TODO: Check if the game is over at the starting position (still needed?)
	game.tetInPlay.Pos.Y++
	err = game.matrix.AddTetrimino(game.tetInPlay)
	if err != nil {
		return nil, fmt.Errorf("failed to add tetrimino to matrix: %w", err)
	}

	return game, nil
}

func (g *Game) MoveLeft() error {
	return g.tetInPlay.MoveLeft(g.matrix)
}

func (g *Game) MoveRight() error {
	return g.tetInPlay.MoveRight(g.matrix)
}

func (g *Game) Rotate(clockwise bool) error {
	return g.tetInPlay.Rotate(g.matrix, clockwise)
}

// Hold will swap the current Tetrimino with the hold Tetrimino.
// If the hold Tetrimino is empty, the current Tetrimino is placed in the hold slot and the addTetInPlay Tetrimino is drawn.
// If not allowed to hold, no action is taken.
// If true is returned the game is over.
func (g *Game) Hold() (bool, error) {
	if !g.canHold {
		return false, nil
	}

	// Remove current Tetrimino
	if err := g.matrix.RemoveTetrimino(g.tetInPlay); err != nil {
		return false, fmt.Errorf("failed to remove hold tetrimino from matrix: %w", err)
	}

	// Swap the current tetrimino with the hold tetrimino
	if g.holdQueue.Value == 0 {
		g.holdQueue = g.tetInPlay
		g.tetInPlay = g.nextQueue.Next()
	} else {
		g.holdQueue, g.tetInPlay = g.tetInPlay, g.holdQueue
	}

	// Add it to the board
	gameOver, err := g.addTetInPlay()
	g.canHold = false // override canHold reset in addTetInPlay()
	if err != nil {
		return false, err
	}
	if gameOver {
		return true, nil
	}

	// Reset the hold tetrimino
	var found bool
	for _, t := range tetris.Tetriminos {
		if t.Value != g.holdQueue.Value {
			continue
		}
		g.holdQueue = t.DeepCopy()
		g.holdQueue.Pos.Y += g.matrix.GetSkyline()
		found = true
		break
	}
	if !found {
		return false, fmt.Errorf("failed to find tetrimino with value '%v'", g.tetInPlay.Value)
	}

	g.canHold = false
	return false, nil
}

// TickLower moves the current Tetrimino down one row. This should be triggered at a regular interval calculated using Fall.
// If the Tetrimino cannot move down, it is locked in place and true is returned.
// Game Over is updated if needed.
func (g *Game) TickLower() (bool, error) {
	lockedDown, err := g.lowerTetInPlay()
	if err != nil {
		return false, fmt.Errorf("failed to lower tetrimino: %w", err)
	}
	if !lockedDown {
		return false, nil
	}

	if g.fall.IsSoftDrop {
		linesCleared := g.tetInPlay.Pos.Y - g.softDropStartRow
		if linesCleared > 0 {
			g.scoring.AddSoftDrop(uint(linesCleared))
		}
	}

	g.tetInPlay = g.nextQueue.Next()
	gameOver, err := g.addTetInPlay()
	if err != nil {
		panic(fmt.Errorf("failed to get addTetInPlay tetrimino (tick): %w", err))
	}

	if gameOver {
		g.gameOver = gameOver
	}

	return gameOver, nil
}

// lowerTetInPlay moves the current Tetrimino down one row if possible.
// If it cannot be moved down this will instead remove completed lines, calculate scores & fall speed and return true (representing a Lock Down).
func (g *Game) lowerTetInPlay() (bool, error) {
	if g.tetInPlay.CanMoveDown(g.matrix) {
		err := g.tetInPlay.MoveDown(&g.matrix)
		if err != nil {
			return false, fmt.Errorf("failed to move tetrimino down: %w", err)
		}
		return false, nil
	}

	action := g.matrix.RemoveCompletedLines(g.tetInPlay)
	if !action.IsValid() {
		return false, fmt.Errorf("invalid action received %q", action.String())
	}
	g.scoring.ProcessAction(action)
	g.fall.CalculateFallSpeeds(g.scoring.Level())

	return true, nil
}

// addTetInPlay will add the tetInPlay to the board at the skyline.
// If possible, it will be moved into the visible matrix (ie. out of the buffer zone).
// If true is returned the game is over.
// canHold is set to true.
func (g *Game) addTetInPlay() (bool, error) {
	// Block Out
	if g.tetInPlay.IsOverlapping(g.matrix) {
		g.gameOver = true
		return true, nil
	}

	if g.tetInPlay.CanMoveDown(g.matrix) {
		g.tetInPlay.Pos.Y++
	} else {
		// Lock Out
		if g.tetInPlay.IsAboveSkyline(g.matrix.GetSkyline()) {
			g.gameOver = true
			return true, nil
		}
	}

	if err := g.matrix.AddTetrimino(g.tetInPlay); err != nil {
		return false, fmt.Errorf("failed to add tetrimino to matrix: %w", err)
	}
	g.canHold = true

	if g.fall.IsSoftDrop {
		g.softDropStartRow = g.tetInPlay.Pos.Y
	}

	return false, nil
}

func (g *Game) HardDrop() (bool, error) {
	g.softDropStartRow = g.tetInPlay.Pos.Y
	for {
		lockedDown, err := g.lowerTetInPlay()
		if err != nil {
			return false, fmt.Errorf("failed to lower tetrimino (hard drop): %w", err)
		}
		if lockedDown {
			break
		}
	}
	linesCleared := g.tetInPlay.Pos.Y - g.softDropStartRow
	if linesCleared > 0 {
		g.scoring.AddHardDrop(uint(g.tetInPlay.Pos.Y - g.softDropStartRow))
	}
	g.softDropStartRow = len(g.matrix)

	g.tetInPlay = g.nextQueue.Next()
	gameOver, err := g.addTetInPlay()
	if err != nil {
		return gameOver, fmt.Errorf("failed to get addTetInPlay tetrimino (hard drop): %w", err)
	}
	return gameOver, nil
}

// ToggleSoftDrop toggles the Soft Drop state of the game.
// If Soft Drop is enabled, the game will calculate the number of lines cleared and add them to the score.
// The time interval for the Fall system is returned.
func (g *Game) ToggleSoftDrop() time.Duration {
	g.fall.ToggleSoftDrop()
	if g.fall.IsSoftDrop {
		g.softDropStartRow = g.tetInPlay.Pos.Y
		return g.fall.SoftDropInterval
	}
	linesCleared := g.tetInPlay.Pos.Y - g.softDropStartRow
	if linesCleared > 0 {
		g.scoring.AddSoftDrop(uint(linesCleared))
	}
	g.softDropStartRow = g.matrix.GetSkyline()
	return g.fall.DefaultInterval
}

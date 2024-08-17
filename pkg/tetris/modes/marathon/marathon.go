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
	ghostTet         *tetris.Tetrimino // The ghost Tetrimino
	holdQueue        *tetris.Tetrimino // The Tetrimino that is being held
	canHold          bool              // Whether the player can hold the current Tetrimino
	gameOver         bool              // Whether the game is over
	softDropStartRow int               // Records where the user began soft drop
	scoring          *tetris.Scoring   // The scoring system
	fall             *tetris.Fall      // The system for calculating the fall speed
}

func NewGame(level, maxLevel uint, gameEnds bool, ghostEnabled bool) (*Game, error) {
	matrix, err := tetris.NewMatrix(40, 10)
	if err != nil {
		return nil, err
	}
	nq := tetris.NewNextQueue(matrix.GetSkyline())

	g := &Game{
		matrix:           matrix,
		nextQueue:        nq,
		tetInPlay:        nq.Next(),
		holdQueue:        tetris.GetEmptyTetrimino(),
		gameOver:         false,
		softDropStartRow: matrix.GetHeight(),
		scoring:          tetris.NewScoring(level, maxLevel, gameEnds),
		fall:             tetris.NewFall(level),
	}

	if ghostEnabled {
		g.ghostTet = g.tetInPlay
	}

	gameOver, err := g.setupNewTetInPlay()
	if err != nil {
		return nil, err
	}
	if gameOver {
		return nil, fmt.Errorf("game over before it began")
	}

	return g, nil
}

func (g *Game) MoveLeft() error {
	_ = g.tetInPlay.MoveLeft(g.matrix)

	err := g.updateGhost()
	if err != nil {
		return fmt.Errorf("failed to update ghost: %w", err)
	}
	return nil
}

func (g *Game) MoveRight() error {
	_ = g.tetInPlay.MoveRight(g.matrix)

	err := g.updateGhost()
	if err != nil {
		return fmt.Errorf("failed to update ghost: %w", err)
	}
	return nil
}

func (g *Game) Rotate(clockwise bool) error {
	err := g.tetInPlay.Rotate(g.matrix, clockwise)
	if err != nil {
		return err
	}

	err = g.updateGhost()
	if err != nil {
		return fmt.Errorf("failed to update ghost: %w", err)
	}
	return nil
}

// Hold will swap the current Tetrimino with the hold Tetrimino.
// If the hold Tetrimino is empty, the current Tetrimino is placed in the hold slot and the setupNewTetInPlay Tetrimino is drawn.
// If not allowed to hold, no action is taken.
// If true is returned the game is over.
func (g *Game) Hold() (bool, error) {
	if !g.canHold {
		return false, nil
	}

	// Swap the current tetrimino with the hold tetrimino
	if g.holdQueue.Value == 0 {
		g.holdQueue = g.tetInPlay
		g.tetInPlay = g.nextQueue.Next()
	} else {
		g.holdQueue, g.tetInPlay = g.tetInPlay, g.holdQueue
	}

	// Add it to the board
	gameOver, err := g.setupNewTetInPlay()
	if err != nil {
		return false, err
	}
	if gameOver {
		return true, nil
	}

	// Reset the hold tetrimino
	t, err := tetris.GetTetrimino(g.holdQueue.Value)
	if err != nil {
		return false, err
	}
	g.holdQueue = t.DeepCopy()
	g.holdQueue.Pos.Y += g.matrix.GetSkyline()

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
	if g.gameOver {
		return true, nil
	}

	if g.fall.IsSoftDrop {
		linesCleared := g.tetInPlay.Pos.Y - g.softDropStartRow
		if linesCleared > 0 {
			g.scoring.AddSoftDrop(uint(linesCleared))
		}
	}

	g.tetInPlay = g.nextQueue.Next()
	gameOver, err := g.setupNewTetInPlay()
	if err != nil {
		panic(fmt.Errorf("failed to get setupNewTetInPlay tetrimino (tick): %w", err))
	}

	if gameOver {
		g.gameOver = gameOver
	}

	return gameOver, nil
}

// lowerTetInPlay moves the current Tetrimino down one row if possible.
// If it cannot be moved down this will instead remove completed lines, calculate scores & fall speed and return true (representing a Lock Down).
// If the max score is reached and the game is configured to end on max level, the Game.gameOver value will be set to true.
func (g *Game) lowerTetInPlay() (bool, error) {
	if g.tetInPlay.MoveDown(g.matrix) {
		return false, nil
	}

	err := g.matrix.AddTetrimino(g.tetInPlay)
	if err != nil {
		return false, err
	}

	action := g.matrix.RemoveCompletedLines(g.tetInPlay)
	if !action.IsValid() {
		return false, fmt.Errorf("invalid action received %q", action.String())
	}

	gameOver := g.scoring.ProcessAction(action)
	if gameOver {
		g.gameOver = true
	}

	g.fall.CalculateFallSpeeds(g.scoring.Level())

	return true, nil
}

// setupNewTetInPlay will do the following setup for the new Tetrimino in play:
//   - If possible, move down one row into the visible Matrix.
//   - Check for Lock Out & Block Out game over conditions.
//   - Reset Game.softDropStartRow if currently Soft Dropping.
//   - Set Game.canHold to true.
//
// It does not modify Game.tetInPlay. If true is returned the game is over.
func (g *Game) setupNewTetInPlay() (bool, error) {
	// Block Out
	if g.tetInPlay.IsOverlapping(g.matrix) {
		g.gameOver = true
		return true, nil
	}

	if !g.tetInPlay.MoveDown(g.matrix) {
		// Lock Out
		if g.tetInPlay.IsAboveSkyline(g.matrix.GetSkyline()) {
			g.gameOver = true
			return true, nil
		}
	}

	g.canHold = true

	if g.fall.IsSoftDrop {
		g.softDropStartRow = g.tetInPlay.Pos.Y
	}

	err := g.updateGhost()
	if err != nil {
		return false, fmt.Errorf("failed to update ghost: %w", err)
	}

	return false, nil
}

func (g *Game) HardDrop() (bool, error) {
	startRow := g.tetInPlay.Pos.Y

	for {
		lockedDown, err := g.lowerTetInPlay()
		if err != nil {
			return false, fmt.Errorf("failed to lower tetrimino (hard drop): %w", err)
		}
		if lockedDown {
			break
		}
	}
	if g.gameOver {
		return true, nil
	}

	linesCleared := g.tetInPlay.Pos.Y - startRow
	g.scoring.AddHardDrop(uint(linesCleared))

	g.tetInPlay = g.nextQueue.Next()
	gameOver, err := g.setupNewTetInPlay()
	if err != nil {
		return gameOver, fmt.Errorf("failed to get setupNewTetInPlay tetrimino (hard drop): %w", err)
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

func (g *Game) updateGhost() error {
	if g.ghostTet == nil {
		return nil
	}

	g.ghostTet = g.tetInPlay.DeepCopy()
	g.ghostTet.Value = 'G'

	for {
		if !g.ghostTet.MoveDown(g.matrix) {
			break
		}
	}

	return nil
}

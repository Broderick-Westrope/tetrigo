package single

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/stuttgart-things/sthings-tetris/pkg/tetris"
)

// Game represents a single player game of Tetris.
// This can be used for Marathon, Sprint, Ultra and other single player modes.
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

type Input struct {
	Level         int  // The starting level of the game.
	MaxLevel      int  // The maximum level the game can reach. 0 means no limit.
	IncreaseLevel bool // Whether the level should increase as the game progresses.
	EndOnMaxLevel bool // Whether the game should end when the maximum level is reached.

	MaxLines      int  // The maximum number of lines to clear before the game ends. 0 means no limit.
	EndOnMaxLines bool // Whether the game should end when the maximum number of lines is cleared.

	GhostEnabled bool       // Whether the ghost Tetrimino should be displayed.
	Rand         *rand.Rand // The random source to use for Tetrimino generation.
}

func NewGame(in *Input) (*Game, error) {
	matrix, err := tetris.NewMatrix(40, 10)
	if err != nil {
		return nil, err
	}
	nq := tetris.NewNextQueue(matrix.GetSkyline(), tetris.WithRandSource(in.Rand))

	scoring, err := tetris.NewScoring(
		in.Level, in.MaxLevel, in.IncreaseLevel, in.EndOnMaxLevel, in.MaxLines, in.EndOnMaxLines,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create scoring system: %w", err)
	}

	g := &Game{
		matrix:           matrix,
		nextQueue:        nq,
		tetInPlay:        nq.Next(),
		holdQueue:        tetris.GetEmptyTetrimino(),
		gameOver:         false,
		softDropStartRow: matrix.GetHeight(),
		scoring:          scoring,
		fall:             tetris.NewFall(in.Level),
	}

	if in.GhostEnabled {
		g.ghostTet = g.tetInPlay
	}

	gameOver := g.setupNewTetInPlay()
	if gameOver {
		return nil, errors.New("game over before it began")
	}

	return g, nil
}

func (g *Game) MoveLeft() {
	_ = g.tetInPlay.MoveLeft(g.matrix)
	g.updateGhost()
}

func (g *Game) MoveRight() {
	_ = g.tetInPlay.MoveRight(g.matrix)
	g.updateGhost()
}

func (g *Game) Rotate(clockwise bool) error {
	err := g.tetInPlay.Rotate(g.matrix, clockwise)
	if err != nil {
		return err
	}

	g.updateGhost()
	return nil
}

// Hold will swap the current Tetrimino with the hold Tetrimino.
// If the hold Tetrimino is empty, the current Tetrimino is placed in the hold slot and
// the setupNewTetInPlay Tetrimino is drawn.
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
	gameOver := g.setupNewTetInPlay()
	if gameOver {
		return true, nil
	}

	// Reset the hold tetrimino
	t, err := tetris.GetTetrimino(g.holdQueue.Value)
	if err != nil {
		return false, err
	}
	g.holdQueue = t.DeepCopy()
	g.holdQueue.Position.Y += g.matrix.GetSkyline()

	g.canHold = false
	return false, nil
}

// TickLower moves the current Tetrimino down one row.
// This should be triggered at a regular interval calculated using Fall.
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
		linesCleared := g.tetInPlay.Position.Y - g.softDropStartRow
		if linesCleared > 0 {
			g.scoring.AddSoftDrop(linesCleared)
		}
	}

	g.tetInPlay = g.nextQueue.Next()
	gameOver := g.setupNewTetInPlay()
	if gameOver {
		g.gameOver = gameOver
	}

	return gameOver, nil
}

func (g *Game) HardDrop() (bool, error) {
	startRow := g.tetInPlay.Position.Y

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

	linesCleared := g.tetInPlay.Position.Y - startRow
	g.scoring.AddHardDrop(linesCleared)

	g.tetInPlay = g.nextQueue.Next()
	gameOver := g.setupNewTetInPlay()
	if gameOver {
		g.gameOver = gameOver
	}

	return gameOver, nil
}

// ToggleSoftDrop toggles the Soft Drop state of the game.
// If Soft Drop is enabled, the game will calculate the number of lines cleared and add them to the score.
func (g *Game) ToggleSoftDrop() {
	g.fall.ToggleSoftDrop()
	if g.fall.IsSoftDrop {
		g.softDropStartRow = g.tetInPlay.Position.Y
		return
	}

	linesCleared := g.tetInPlay.Position.Y - g.softDropStartRow
	if linesCleared > 0 {
		g.scoring.AddSoftDrop(linesCleared)
	}
	g.softDropStartRow = g.matrix.GetSkyline()
}

// GetFallInterval returns the time interval for the Fall system.
func (g *Game) GetFallInterval() time.Duration {
	if g.fall.IsSoftDrop {
		return g.fall.SoftDropInterval
	}
	return g.fall.DefaultInterval
}

// EndGame sets Game.gameOver to true.
func (g *Game) EndGame() {
	g.gameOver = true
}

// lowerTetInPlay moves the current Tetrimino down one row if possible.
// If it cannot be moved down this will instead remove completed lines, calculate scores and
// fall speed and return true (representing a Lock Down).
// If the max score is reached and the game is configured to end on max level,
// the Game.gameOver value will be set to true.
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

	gameOver, err := g.scoring.ProcessAction(action)
	if err != nil {
		return false, fmt.Errorf("failed to process action: %w", err)
	}
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
func (g *Game) setupNewTetInPlay() bool {
	// Block Out
	if !g.tetInPlay.IsValid(g.matrix, false) {
		g.gameOver = true
		return true
	}

	if !g.tetInPlay.MoveDown(g.matrix) {
		// Lock Out
		if g.tetInPlay.IsAboveSkyline(g.matrix.GetSkyline()) {
			g.gameOver = true
			return true
		}
	}

	g.canHold = true

	if g.fall.IsSoftDrop {
		g.softDropStartRow = g.tetInPlay.Position.Y
	}

	g.updateGhost()
	return false
}

func (g *Game) updateGhost() {
	if g.ghostTet == nil {
		return
	}

	g.ghostTet = g.tetInPlay.DeepCopy()
	g.ghostTet.Value = 'G'

	for {
		if !g.ghostTet.MoveDown(g.matrix) {
			break
		}
	}
}

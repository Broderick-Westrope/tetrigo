package marathon

import (
	"fmt"
	"time"

	"github.com/Broderick-Westrope/tetrigo/pkg/tetris"
)

type Game struct {
	matrix     tetris.Matrix     // The Matrix of cells on which the game is played
	bag        *tetris.Bag       // The bag of upcoming Tetriminos
	currentTet *tetris.Tetrimino // The current Tetrimino in play
	holdTet    *tetris.Tetrimino // The Tetrimino that is being held
	canHold    bool              // Whether the player can hold the current Tetrimino
	gameOver   bool              // Whether the game is over
	startRow   int               // The line at which the game started
	scoring    *tetris.Scoring   // The scoring system
	fall       *tetris.Fall      // The system for calculating the fall speed
}

func NewGame(level, maxLevel uint) (*Game, error) {
	matrix, err := tetris.NewMatrix(40, 10)
	if err != nil {
		return nil, err
	}
	bag := tetris.NewBag(matrix.GetSkyline())

	game := &Game{
		matrix:     matrix,
		bag:        bag,
		currentTet: bag.Next(),
		holdTet:    tetris.EmptyTetrimino,
		canHold:    true,
		gameOver:   false,
		// TODO: is start line needed?
		startRow: matrix.GetHeight(),
		scoring:  tetris.NewScoring(level, maxLevel),
		fall:     tetris.NewFall(level),
	}

	// TODO: Check if the game is over at the starting position (still needed?)
	game.currentTet.Pos.Y++
	err = game.matrix.AddTetrimino(game.currentTet)
	if err != nil {
		return nil, fmt.Errorf("failed to add tetrimino to matrix: %w", err)
	}

	return game, nil
}

func (g *Game) MoveLeft() error {
	return g.currentTet.MoveLeft(g.matrix)
}

func (g *Game) MoveRight() error {
	return g.currentTet.MoveRight(g.matrix)
}

func (g *Game) Rotate(clockwise bool) error {
	return g.currentTet.Rotate(g.matrix, clockwise)
}

// Hold will swap the current Tetrimino with the hold Tetrimino.
// If the hold Tetrimino is empty, the current Tetrimino is placed in the hold slot and the next Tetrimino is drawn.
// If not allowed to hold, no action is taken.
func (g *Game) Hold() error {
	if !g.canHold {
		return nil
	}

	// Swap the current tetrimino with the hold tetrimino
	if g.holdTet.Value == 0 {
		g.holdTet = g.currentTet
		g.currentTet = g.bag.Next()
	} else {
		g.holdTet, g.currentTet = g.currentTet, g.holdTet
	}

	if err := g.matrix.RemoveTetrimino(g.holdTet); err != nil {
		return fmt.Errorf("failed to remove hold tetrimino from matrix: %w", err)
	}

	// Reset the position of the hold tetrimino
	var found bool
	for _, t := range tetris.Tetriminos {
		if t.Value != g.holdTet.Value {
			continue
		}
		g.holdTet.Pos = t.Pos
		g.holdTet.Pos.Y += g.matrix.GetSkyline()
		found = true
		break
	}
	if !found {
		return fmt.Errorf("failed to find tetrimino with value '%v'", g.currentTet.Value)
	}

	// Add the current tetrimino to the matrix
	err := g.matrix.AddTetrimino(g.currentTet)
	if err != nil {
		return fmt.Errorf("failed to add tetrimino to matrix: %w", err)
	}

	g.canHold = false
	return nil
}

// TickLower moves the current Tetrimino down one row. This should be triggered at a regular interval calculated using Fall.
// If the Tetrimino cannot move down, it is locked in place and true is returned.
// Game Over is updated if the Tetrimino is blocked out.
func (g *Game) TickLower() (bool, error) {
	lockedOut, err := g.lowerCurrentTet()
	if err != nil {
		return false, fmt.Errorf("failed to lower tetrimino: %w", err)
	}

	if !lockedOut {
		return false, nil
	}

	if g.fall.IsSoftDrop {
		linesCleared := g.currentTet.Pos.Y - g.startRow
		if linesCleared > 0 {
			g.scoring.AddSoftDrop(uint(linesCleared))
		}
	}

	gameOver, err := g.next()
	if err != nil {
		panic(fmt.Errorf("failed to get next tetrimino (tick): %w", err))
	}

	if gameOver {
		g.gameOver = gameOver
	}

	return true, nil
}

func (g *Game) lowerCurrentTet() (bool, error) {
	if g.currentTet.CanMoveDown(g.matrix) {
		err := g.currentTet.MoveDown(&g.matrix)
		if err != nil {
			return false, fmt.Errorf("failed to move tetrimino down: %w", err)
		}
		return false, nil
	}

	action := g.matrix.RemoveCompletedLines(g.currentTet)
	g.scoring.ProcessAction(action)
	g.fall.CalculateFallSpeeds(g.scoring.Level())

	return true, nil
}

func (g *Game) next() (bool, error) {
	g.currentTet = g.bag.Next()

	// Block Out
	if g.currentTet.IsOverlapping(g.matrix) {
		return true, nil
	}

	if g.currentTet.CanMoveDown(g.matrix) {
		g.currentTet.Pos.Y++
	} else {
		// Lock Out
		if g.currentTet.IsAbovePlayfield(g.matrix.GetSkyline()) {
			return true, nil
		}
	}

	if err := g.matrix.AddTetrimino(g.currentTet); err != nil {
		return false, fmt.Errorf("failed to add tetrimino to matrix: %w", err)
	}
	g.canHold = true

	if g.fall.IsSoftDrop {
		g.startRow = g.currentTet.Pos.Y
	}

	return false, nil
}

func (g *Game) HardDrop() (bool, error) {
	g.startRow = g.currentTet.Pos.Y
	for {
		lockedDown, err := g.lowerCurrentTet()
		if err != nil {
			return false, fmt.Errorf("failed to lower tetrimino (hard drop): %w", err)
		}
		if lockedDown {
			break
		}
	}
	linesCleared := g.currentTet.Pos.Y - g.startRow
	if linesCleared > 0 {
		g.scoring.AddHardDrop(uint(g.currentTet.Pos.Y - g.startRow))
	}
	g.startRow = len(g.matrix)

	gameOver, err := g.next()
	if err != nil {
		return gameOver, fmt.Errorf("failed to get next tetrimino (hard drop): %w", err)
	}
	return gameOver, nil
}

// ToggleSoftDrop toggles the Soft Drop state of the game.
// If Soft Drop is enabled, the game will calculate the number of lines cleared and add them to the score.
// The time interval for the Fall system is returned.
func (g *Game) ToggleSoftDrop() time.Duration {
	g.fall.ToggleSoftDrop()
	if g.fall.IsSoftDrop {
		g.startRow = g.currentTet.Pos.Y
		return g.fall.SoftDropInterval
	}
	linesCleared := g.currentTet.Pos.Y - g.startRow
	if linesCleared > 0 {
		g.scoring.AddSoftDrop(uint(linesCleared))
	}
	g.startRow = g.matrix.GetSkyline()
	return g.fall.DefaultInterval
}

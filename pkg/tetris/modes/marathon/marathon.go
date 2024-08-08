package marathon

import (
	"fmt"
	"time"

	"github.com/Broderick-Westrope/tetrigo/pkg/tetris"
)

type Game struct {
	Matrix     tetris.Matrix     // The Matrix of cells on which the game is played
	Bag        *tetris.Bag       // The Bag of upcoming Tetriminos
	CurrentTet *tetris.Tetrimino // The current Tetrimino in play
	HoldTet    *tetris.Tetrimino // The Tetrimino that is being held
	canHold    bool              // Whether the player can hold the current Tetrimino
	gameOver   bool              // Whether the game is over
	maxLevel   uint              // The maximum level that can be reached
	StartRow   int               // The line at which the game started
	Scoring    *tetris.Scoring   // The Scoring system
	Fall       *tetris.Fall      // The system for calculating the Fall speed
}

func NewGame(level, maxLevel uint) (*Game, error) {
	matrix, err := tetris.NewMatrix(40, 10)
	if err != nil {
		return nil, err
	}
	bag := tetris.NewBag(matrix.GetStartRow())

	game := &Game{
		Matrix:     matrix,
		Bag:        bag,
		CurrentTet: bag.Next(),
		HoldTet:    tetris.EmptyTetrimino,
		canHold:    true,
		gameOver:   false,
		maxLevel:   maxLevel,
		// TODO: is start line needed?
		StartRow: matrix.GetHeight(),
		Scoring:  tetris.NewScoring(level),
		Fall:     tetris.NewFall(level),
	}

	// TODO: Check if the game is over at the starting position (still needed?)
	game.CurrentTet.Pos.Y++
	err = game.Matrix.AddTetrimino(game.CurrentTet)
	if err != nil {
		return nil, fmt.Errorf("failed to add tetrimino to Matrix: %w", err)
	}

	return game, nil
}

func (g *Game) IsGameOver() bool {
	return g.gameOver
}

func (g *Game) MoveLeft() error {
	return g.CurrentTet.MoveLeft(g.Matrix)
}

func (g *Game) MoveRight() error {
	return g.CurrentTet.MoveRight(g.Matrix)
}

func (g *Game) Rotate(clockwise bool) error {
	return g.CurrentTet.Rotate(g.Matrix, clockwise)
}

func (g *Game) Hold() error {
	if !g.canHold {
		return nil
	}

	// Swap the current tetrimino with the hold tetrimino
	if g.HoldTet.Value == 0 {
		g.HoldTet = g.CurrentTet
		g.CurrentTet = g.Bag.Next()
	} else {
		g.HoldTet, g.CurrentTet = g.CurrentTet, g.HoldTet
	}

	if err := g.Matrix.RemoveTetrimino(g.HoldTet); err != nil {
		return fmt.Errorf("failed to remove hold tetrimino from matrix: %w", err)
	}

	// Reset the position of the hold tetrimino
	var found bool
	for _, t := range tetris.Tetriminos {
		if t.Value != g.HoldTet.Value {
			continue
		}
		g.HoldTet.Pos = t.Pos
		g.HoldTet.Pos.Y += g.Matrix.GetStartRow()
		found = true
		break
	}
	if !found {
		return fmt.Errorf("failed to find tetrimino with value '%v'", g.CurrentTet.Value)
	}

	// Add the current tetrimino to the matrix
	err := g.Matrix.AddTetrimino(g.CurrentTet)
	if err != nil {
		return fmt.Errorf("failed to add tetrimino to matrix: %w", err)
	}

	g.canHold = false
	return nil
}

func (g *Game) Lower() (bool, error) {
	if g.CurrentTet.CanMoveDown(g.Matrix) {
		err := g.CurrentTet.MoveDown(&g.Matrix)
		if err != nil {
			return false, fmt.Errorf("failed to move tetrimino down: %w", err)
		}
		return false, nil
	}

	action := g.Matrix.RemoveCompletedLines(g.CurrentTet)
	g.Scoring.ProcessAction(action, g.maxLevel)

	g.Fall.CalculateFallSpeeds(g.Scoring.Level())

	if g.Fall.IsSoftDrop {
		linesCleared := g.CurrentTet.Pos.Y - g.StartRow
		if linesCleared > 0 {
			g.Scoring.AddSoftDrop(uint(linesCleared))
		}
	}

	gameOver, err := g.Next()
	if err != nil {
		panic(fmt.Errorf("failed to get next tetrimino (tick): %w", err))
	}

	if gameOver {
		g.gameOver = gameOver
	}

	return true, nil
}

func (g *Game) Next() (bool, error) {
	g.CurrentTet = g.Bag.Next()

	// Block Out
	if g.CurrentTet.IsOverlapping(&g.Matrix) {
		return true, nil
	}

	if g.CurrentTet.CanMoveDown(g.Matrix) {
		g.CurrentTet.Pos.Y++
	} else {
		// Lock Out
		if g.CurrentTet.IsAbovePlayfield(len(g.Matrix)) {
			return true, nil
		}
	}

	if err := g.Matrix.AddTetrimino(g.CurrentTet); err != nil {
		return false, fmt.Errorf("failed to add tetrimino to matrix: %w", err)
	}
	g.canHold = true

	if g.Fall.IsSoftDrop {
		g.StartRow = g.CurrentTet.Pos.Y
	}

	return false, nil
}

func (g *Game) HardDrop() (bool, error) {
	g.StartRow = g.CurrentTet.Pos.Y
	for {
		lockedDown, err := g.Lower()
		if err != nil {
			return false, fmt.Errorf("failed to lower tetrimino (hard drop): %w", err)
		}
		if lockedDown {
			break
		}
	}
	linesCleared := g.CurrentTet.Pos.Y - g.StartRow
	if linesCleared > 0 {
		g.Scoring.AddHardDrop(uint(g.CurrentTet.Pos.Y - g.StartRow))
	}
	g.StartRow = len(g.Matrix)

	gameOver, err := g.Next()
	if err != nil {
		return gameOver, fmt.Errorf("failed to get next tetrimino (hard drop): %w", err)
	}
	return gameOver, nil
}

func (g *Game) ToggleSoftDrop() time.Duration {
	g.Fall.ToggleSoftDrop()
	if g.Fall.IsSoftDrop {
		g.StartRow = g.CurrentTet.Pos.Y
		return g.Fall.SoftDropTime
	}
	linesCleared := g.CurrentTet.Pos.Y - g.StartRow
	if linesCleared > 0 {
		g.Scoring.AddSoftDrop(uint(linesCleared))
	}
	g.StartRow = g.Matrix.GetStartRow()
	return g.Fall.DefaultTime
}

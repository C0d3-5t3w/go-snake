package game

import (
	"math/rand"
	"time"

	"github.com/C0d3-5t3w/go-snake/internal/config"
)

// Direction represents movement direction in 3D space
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
	Forward
	Backward
)

// Point3D represents a position in 3D space
type Point3D struct {
	X, Y, Z int
}

// Snake represents the player's snake
type Snake struct {
	Body      []Point3D
	Direction Direction
	GrowCount int
}

// GameState represents the current state of the game
type GameState int

const (
	Playing GameState = iota
	Paused
	GameOver
)

// Game represents the snake game
type Game struct {
	Config        *config.Config
	Snake         Snake
	Food          Point3D
	Grid          int
	Score         int
	State         GameState
	Speed         float64
	LastUpdate    time.Time
	OnScoreChange func(int)
}

// NewGame creates a new game instance
func NewGame(cfg *config.Config) *Game {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		Config: cfg,
		Grid:   cfg.Game.GridSize,
		Speed:  cfg.Game.InitialSpeed,
		State:  Paused,
	}

	game.Reset()
	return game
}

// Reset resets the game to initial state
func (g *Game) Reset() {
	// Create snake in the center of the grid
	center := g.Grid / 2
	g.Snake = Snake{
		Body: []Point3D{
			{X: center, Y: center, Z: center},
		},
		Direction: Forward,
	}

	// Grow snake to initial length
	g.Snake.GrowCount = g.Config.Game.InitialLength - 1

	// Place food
	g.PlaceFood()

	// Reset score and speed
	g.Score = 0
	g.Speed = g.Config.Game.InitialSpeed
	g.State = Playing
	g.LastUpdate = time.Now()

	// Notify score change
	if g.OnScoreChange != nil {
		g.OnScoreChange(g.Score)
	}
}

// PlaceFood places food at a random position not occupied by the snake
func (g *Game) PlaceFood() {
	for {
		// Generate random position
		food := Point3D{
			X: rand.Intn(g.Grid),
			Y: rand.Intn(g.Grid),
			Z: rand.Intn(g.Grid),
		}

		// Check if position overlaps with snake
		overlap := false
		for _, part := range g.Snake.Body {
			if part.X == food.X && part.Y == food.Y && part.Z == food.Z {
				overlap = true
				break
			}
		}

		if !overlap {
			g.Food = food
			break
		}
	}
}

// ChangeDirection changes the snake's direction
func (g *Game) ChangeDirection(dir Direction) {
	// Prevent 180-degree turns
	opposites := map[Direction]Direction{
		Up:       Down,
		Down:     Up,
		Left:     Right,
		Right:    Left,
		Forward:  Backward,
		Backward: Forward,
	}

	if opposites[dir] != g.Snake.Direction {
		g.Snake.Direction = dir
	}
}

// Update updates the game state
func (g *Game) Update() bool {
	if g.State != Playing {
		return false
	}

	now := time.Now()
	updateInterval := time.Duration(1000/g.Speed) * time.Millisecond

	if now.Sub(g.LastUpdate) < updateInterval {
		return false
	}

	g.LastUpdate = now

	// Save head position before moving
	head := g.Snake.Body[0]

	// Calculate new head position
	newHead := Point3D{X: head.X, Y: head.Y, Z: head.Z}

	// Move based on direction
	switch g.Snake.Direction {
	case Up:
		newHead.Y++
	case Down:
		newHead.Y--
	case Left:
		newHead.X--
	case Right:
		newHead.X++
	case Forward:
		newHead.Z++
	case Backward:
		newHead.Z--
	}

	// Check for wall collision
	if newHead.X < 0 || newHead.X >= g.Grid ||
		newHead.Y < 0 || newHead.Y >= g.Grid ||
		newHead.Z < 0 || newHead.Z >= g.Grid {
		g.State = GameOver
		return true
	}

	// Check for self collision
	for _, part := range g.Snake.Body {
		if newHead.X == part.X && newHead.Y == part.Y && newHead.Z == part.Z {
			g.State = GameOver
			return true
		}
	}

	// Check for food collision
	ateFood := newHead.X == g.Food.X && newHead.Y == g.Food.Y && newHead.Z == g.Food.Z

	// Add new head to the snake
	g.Snake.Body = append([]Point3D{newHead}, g.Snake.Body...)

	// If food was eaten or snake is still growing
	if ateFood || g.Snake.GrowCount > 0 {
		if ateFood {
			g.Score += 10
			g.PlaceFood()
			g.Snake.GrowCount++

			// Increase speed
			if g.Speed < g.Config.Game.MaxSpeed {
				g.Speed += g.Config.Game.SpeedIncrement
			}

			// Notify score change
			if g.OnScoreChange != nil {
				g.OnScoreChange(g.Score)
			}
		}

		if g.Snake.GrowCount > 0 {
			g.Snake.GrowCount--
		}
	} else {
		// Remove tail if not growing
		g.Snake.Body = g.Snake.Body[:len(g.Snake.Body)-1]
	}

	return true
}

// TogglePause toggles the pause state
func (g *Game) TogglePause() {
	if g.State == Playing {
		g.State = Paused
	} else if g.State == Paused {
		g.State = Playing
		g.LastUpdate = time.Now()
	}
}

// IsGameOver returns true if the game is over
func (g *Game) IsGameOver() bool {
	return g.State == GameOver
}

package gui

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"

	"github.com/C0d3-5t3w/go-snake/internal/config"
	"github.com/C0d3-5t3w/go-snake/internal/game"
	"github.com/C0d3-5t3w/go-snake/internal/storage"
)

const (
	tileSize = 20 // Size of each grid tile in pixels
)

// EbitenGame holds the game state for Ebitengine
type EbitenGame struct {
	game      *game.Game
	config    *config.Config
	storage   *storage.Storage
	tileSize  int
	infoFont  font.Face
	lastFrame time.Time

	// Cached images for performance
	snakeHeadImg *ebiten.Image
	snakeBodyImg *ebiten.Image
	foodImg      *ebiten.Image
	bgImg        *ebiten.Image
	gridImg      *ebiten.Image // Optional grid image
}

// NewEbitenGUI initializes the Ebiten game wrapper
func NewEbitenGUI(g *game.Game, cfg *config.Config, s *storage.Storage) (*EbitenGame, error) {
	// Basic font for drawing text
	infoFont := basicfont.Face7x13

	eg := &EbitenGame{
		game:     g,
		config:   cfg,
		storage:  s,
		tileSize: tileSize,
		infoFont: infoFont,
	}

	// Initialize Ebiten window settings
	ebiten.SetWindowSize(cfg.Graphics.WindowWidth, cfg.Graphics.WindowHeight)
	ebiten.SetWindowTitle("Go Snake 2D")
	ebiten.SetVsyncEnabled(cfg.Graphics.Vsync)
	if cfg.Graphics.Fullscreen {
		ebiten.SetFullscreen(true)
	}

	// Pre-render images for drawing elements
	eg.createImages()

	// Set game callbacks (if needed, e.g., score updates)
	g.OnScoreChange = func(score int) {
		// Score is drawn directly in the Draw method, no need for label update
	}

	return eg, nil
}

// Run starts the Ebitengine game loop
func (eg *EbitenGame) Run() error {
	return ebiten.RunGame(eg)
}

// createImages pre-renders simple images for snake, food, etc.
func (eg *EbitenGame) createImages() {
	s := eg.tileSize

	// Snake Head (using config color)
	hc := eg.config.Colors.SnakeHead
	eg.snakeHeadImg = ebiten.NewImage(s, s)
	vector.DrawFilledRect(eg.snakeHeadImg, 0, 0, float32(s), float32(s), color.RGBA{R: uint8(hc[0] * 255), G: uint8(hc[1] * 255), B: uint8(hc[2] * 255), A: 255}, false)

	// Snake Body (using config color)
	bc := eg.config.Colors.SnakeBody
	eg.snakeBodyImg = ebiten.NewImage(s, s)
	vector.DrawFilledRect(eg.snakeBodyImg, 0, 0, float32(s), float32(s), color.RGBA{R: uint8(bc[0] * 255), G: uint8(bc[1] * 255), B: uint8(bc[2] * 255), A: 255}, false)

	// Food (using config color)
	fc := eg.config.Colors.Food
	eg.foodImg = ebiten.NewImage(s, s)
	vector.DrawFilledRect(eg.foodImg, 0, 0, float32(s), float32(s), color.RGBA{R: uint8(fc[0] * 255), G: uint8(fc[1] * 255), B: uint8(fc[2] * 255), A: 255}, false)

	// Background (using config color)
	bgc := eg.config.Colors.Background
	eg.bgImg = ebiten.NewImage(1, 1) // Create a 1x1 pixel image for the background color
	eg.bgImg.Fill(color.RGBA{R: uint8(bgc[0] * 255), G: uint8(bgc[1] * 255), B: uint8(bgc[2] * 255), A: 255})

	// Grid (optional, draw lines)
	// We can draw the grid dynamically in the Draw function instead for simplicity
}

// Update proceeds the game state.
func (eg *EbitenGame) Update() error {
	// Handle input
	eg.handleInput()

	// Update game logic
	// Note: Ebiten's Update function is called 60 times per second by default.
	// The internal game Update() has its own speed control based on LastUpdate.
	// This is fine, the game logic will only advance when its internal timer allows.
	eg.game.Update()

	return nil
}

// handleInput processes user input
func (eg *EbitenGame) handleInput() {
	// Game controls
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		eg.game.TogglePause()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		if eg.game.IsGameOver() {
			eg.game.Reset()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Note: Ebiten handles closing the window, maybe map this to pause or menu?
		log.Println("Escape pressed - exiting game (Ebiten handles window close)")
		// Or, implement a quit confirm dialog later
	}

	// Movement controls - only process if playing
	if eg.game.State == game.Playing {
		// Map W to Up (Y-)
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			eg.game.ChangeDirection(game.Up)
		}
		// Map S to Down (Y+)
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			eg.game.ChangeDirection(game.Down)
		}
		// Map A to Left (X-)
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			eg.game.ChangeDirection(game.Left)
		}
		// Map D to Right (X+)
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			eg.game.ChangeDirection(game.Right)
		}
	}
}

// Draw draws the game screen.
func (eg *EbitenGame) Draw(screen *ebiten.Image) {
	// Draw background
	screenW, screenH := screen.Size()
	bgOpts := &ebiten.DrawImageOptions{}
	bgOpts.GeoM.Scale(float64(screenW), float64(screenH)) // Scale the 1x1 pixel bg image
	screen.DrawImage(eg.bgImg, bgOpts)

	// Calculate offsets to center the grid
	gridWidth := eg.game.Grid * eg.tileSize
	gridHeight := eg.game.Grid * eg.tileSize
	offsetX := (screenW - gridWidth) / 2
	offsetY := (screenH - gridHeight) / 2

	// Draw grid lines (optional, can be resource intensive)
	gridC := eg.config.Colors.Grid
	gridColor := color.RGBA{R: uint8(gridC[0] * 255), G: uint8(gridC[1] * 255), B: uint8(gridC[2] * 255), A: 100} // Semi-transparent grid
	for i := 0; i <= eg.game.Grid; i++ {
		fx := float32(offsetX + i*eg.tileSize)
		fy := float32(offsetY + i*eg.tileSize)
		// Vertical line
		vector.StrokeLine(screen, fx, float32(offsetY), fx, float32(offsetY+gridHeight), 1, gridColor, false)
		// Horizontal line
		vector.StrokeLine(screen, float32(offsetX), fy, float32(offsetX+gridWidth), fy, 1, gridColor, false)
	}

	// Draw snake
	for i, part := range eg.game.Snake.Body {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(offsetX+part.X*eg.tileSize), float64(offsetY+part.Y*eg.tileSize))
		if i == 0 {
			screen.DrawImage(eg.snakeHeadImg, opts)
		} else {
			screen.DrawImage(eg.snakeBodyImg, opts)
		}
	}

	// Draw food
	foodOpts := &ebiten.DrawImageOptions{}
	foodOpts.GeoM.Translate(float64(offsetX+eg.game.Food.X*eg.tileSize), float64(offsetY+eg.game.Food.Y*eg.tileSize))
	screen.DrawImage(eg.foodImg, foodOpts)

	// Draw score and status
	scoreText := fmt.Sprintf("Score: %d", eg.game.Score)
	statusText := ""
	switch eg.game.State {
	case game.Playing:
		statusText = "Playing - WASD: Move"
	case game.Paused:
		statusText = "Paused - Press P to Start"
	case game.GameOver:
		statusText = fmt.Sprintf("Game Over - Score: %d - Press R", eg.game.Score)
	}

	// Draw text using ebitenutil for simplicity
	ebitenutil.DebugPrintAt(screen, scoreText, 10, 10)
	ebitenutil.DebugPrintAt(screen, statusText, 10, 30)

	// Draw FPS counter
	fps := ebiten.ActualFPS()
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.1f", fps), screenW-100, 10)
}

// Layout takes the outside size (e.g., window size) and returns the (logical) screen size.
func (eg *EbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Use the configured window size as the logical size
	return eg.config.Graphics.WindowWidth, eg.config.Graphics.WindowHeight
}

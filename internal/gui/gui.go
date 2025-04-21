package gui

import (
	"fmt"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"

	"github.com/C0d3-5t3w/go-snake/internal/config"
	"github.com/C0d3-5t3w/go-snake/internal/game"
	"github.com/C0d3-5t3w/go-snake/internal/storage"
)

// GUI represents the graphical user interface
type GUI struct {
	app          *app.Application
	scene        *core.Node
	camera       *camera.Camera
	game         *game.Game
	config       *config.Config
	storage      *storage.Storage
	snakeMeshes  map[int]*graphic.Mesh
	foodMesh     *graphic.Mesh
	scoreLabel   *gui.Label
	statusLabel  *gui.Label
	orbitControl *camera.OrbitControl
	cameraTarget *core.Node
	gridHelpers  []*helper.Axes
}

// NewGUI creates a new GUI instance
func NewGUI(g *game.Game, cfg *config.Config, s *storage.Storage) (*GUI, error) {
	// Create application
	a := app.App()

	// Set application options
	a.IWindow.(*window.GlfwWindow).SetSize(cfg.Graphics.WindowWidth, cfg.Graphics.WindowHeight)
	if cfg.Graphics.Fullscreen {
		a.IWindow.(*window.GlfwWindow).SetFullscreen(true)
	}

	// Create scene
	scene := core.NewNode()

	// Create camera
	cam := camera.New(1)
	cam.SetPosition(0, 15, -20)
	scene.Add(cam)

	// Create camera target (for orbit control)
	target := core.NewNode()
	scene.Add(target)

	// Create orbit control
	orbitControl := camera.NewOrbitControl(cam)

	// Create GUI
	gui.Manager().Set(scene)

	// Create score label
	scoreLabel := gui.NewLabel("Score: 0")
	scoreLabel.SetPosition(10, 10)
	scoreLabel.SetColor(math32.NewColor("white"))
	scoreLabel.SetFontSize(20)

	// Create status label
	statusLabel := gui.NewLabel("Game Paused - Press P to Start")
	statusLabel.SetPosition(float32(cfg.Graphics.WindowWidth/2-100), 10)
	statusLabel.SetColor(math32.NewColor("white"))
	statusLabel.SetFontSize(20)

	// Create and add lights
	ambientLight := light.NewAmbient(&math32.Color{R: 0.4, G: 0.4, B: 0.4}, 1.0)
	scene.Add(ambientLight)

	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(float32(g.Grid)*1.5, float32(g.Grid)*2, float32(g.Grid)*1.5)
	scene.Add(pointLight)

	dirLight := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 1.0)
	dirLight.SetPosition(float32(g.Grid), float32(g.Grid), float32(g.Grid))
	scene.Add(dirLight)

	// Create grid helpers
	gridHelpers := make([]*helper.Axes, 0)
	for i := 0; i < g.Grid+1; i++ {
		for j := 0; j < g.Grid+1; j++ {
			axesHelper := helper.NewAxes(0.1)
			axesHelper.SetPosition(float32(i), 0, float32(j))
			scene.Add(axesHelper)
			gridHelpers = append(gridHelpers, axesHelper)
		}
	}

	gui := &GUI{
		app:          a,
		scene:        scene,
		camera:       cam,
		game:         g,
		config:       cfg,
		storage:      s,
		snakeMeshes:  make(map[int]*graphic.Mesh),
		scoreLabel:   scoreLabel,
		statusLabel:  statusLabel,
		orbitControl: orbitControl,
		cameraTarget: target,
		gridHelpers:  gridHelpers,
	}

	// Create food mesh
	gui.createFoodMesh()

	// Set game callback for score changes
	g.OnScoreChange = func(score int) {
		scoreLabel.SetText(fmt.Sprintf("Score: %d", score))
	}

	return gui, nil
}

// createFoodMesh creates the 3D mesh for food
func (g *GUI) createFoodMesh() {
	// Create sphere for food
	geom := geometry.NewSphere(0.4, 16, 16)
	mat := material.NewStandard(&math32.Color{
		R: g.config.Colors.Food[0],
		G: g.config.Colors.Food[1],
		B: g.config.Colors.Food[2],
	})

	// Add glow effect
	mat.SetOpacity(0.9)
	mat.SetEmissiveColor(&math32.Color{R: 0.5, G: 0, B: 0})
	mat.SetSpecularColor(&math32.Color{R: 1, G: 0.3, B: 0.3})
	mat.SetShininess(50)

	mesh := graphic.NewMesh(geom, mat)

	// Position at food location
	mesh.SetPosition(
		float32(g.game.Food.X),
		float32(g.game.Food.Y),
		float32(g.game.Food.Z),
	)

	g.scene.Add(mesh)
	g.foodMesh = mesh
}

// updateFoodPosition updates the position of the food mesh
func (g *GUI) updateFoodPosition() {
	if g.foodMesh != nil {
		g.scene.Remove(g.foodMesh)
	}
	g.createFoodMesh()
}

// updateSnakeMeshes updates the snake's meshes
func (g *GUI) updateSnakeMeshes() {
	// Clear old meshes
	for _, mesh := range g.snakeMeshes {
		g.scene.Remove(mesh)
	}
	g.snakeMeshes = make(map[int]*graphic.Mesh)

	// Create new meshes for each snake segment
	for i, part := range g.game.Snake.Body {
		var geom geometry.IGeometry
		var mat material.IMaterial

		if i == 0 {
			// Head is a slightly larger cube with eyes
			geom = geometry.NewBox(0.9, 0.9, 0.9)
			mat = material.NewStandard(&math32.Color{
				R: g.config.Colors.SnakeHead[0],
				G: g.config.Colors.SnakeHead[1],
				B: g.config.Colors.SnakeHead[2],
			})
			mat.(*material.Standard).SetSpecularColor(&math32.Color{R: 0.5, G: 1, B: 0.5})
			mat.(*material.Standard).SetShininess(30)
		} else {
			// Body is a regular cube
			geom = geometry.NewBox(0.8, 0.8, 0.8)
			mat = material.NewStandard(&math32.Color{
				R: g.config.Colors.SnakeBody[0],
				G: g.config.Colors.SnakeBody[1],
				B: g.config.Colors.SnakeBody[2],
			})
		}

		mesh := graphic.NewMesh(geom, mat)
		mesh.SetPosition(
			float32(part.X),
			float32(part.Y),
			float32(part.Z),
		)

		g.scene.Add(mesh)
		g.snakeMeshes[i] = mesh
	}
}

// Run starts the GUI main loop
func (g *GUI) Run() {
	// Handle key events
	g.app.Subscribe(window.OnKeyDown, g.onKeyDown)

	// Combine update and render in the Run callback
	g.app.Run(func(rend *renderer.Renderer, deltaTime time.Duration) {
		g.game.Update()
		g.updateSnakeMeshes()
		g.updateFoodPosition()

		// Update status label text
		switch g.game.State {
		case game.Playing:
			g.statusLabel.SetText("Playing")
		case game.Paused:
			g.statusLabel.SetText("Game Paused - Press P to Start")
		case game.GameOver:
			g.statusLabel.SetText(fmt.Sprintf("Game Over - Score: %d - Press R to Restart", g.game.Score))
			g.storage.AddHighScore("Player", g.game.Score)
			g.storage.Save()
		}

		// Camera follow
		head := g.game.Snake.Body[0]
		targetPos := math32.NewVector3(float32(head.X), float32(head.Y), float32(head.Z))
		g.orbitControl.SetTarget(*targetPos)

		// Render scene (background/clear handled internally)
		rend.Render(g.scene, g.camera)
	})
}

// onKeyDown handles key press events
func (g *GUI) onKeyDown(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)

	// Game controls
	switch kev.Key {
	case window.KeyP:
		g.game.TogglePause()
	case window.KeyR:
		if g.game.IsGameOver() {
			g.game.Reset()
		}
	case window.KeyEscape:
		g.app.Exit()
	}

	// Movement controls
	if g.game.State == game.Playing {
		switch {
		case kev.Key == window.KeyW:
			g.game.ChangeDirection(game.Forward)
		case kev.Key == window.KeyS:
			g.game.ChangeDirection(game.Backward)
		case kev.Key == window.KeyA:
			g.game.ChangeDirection(game.Left)
		case kev.Key == window.KeyD:
			g.game.ChangeDirection(game.Right)
		case kev.Key == window.KeySpace:
			g.game.ChangeDirection(game.Up)
		case kev.Key == window.KeyLeftShift:
			g.game.ChangeDirection(game.Down)
		}
	}
}

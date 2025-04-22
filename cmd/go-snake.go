package main

import (
	"log"

	"github.com/C0d3-5t3w/go-snake/internal/config"
	"github.com/C0d3-5t3w/go-snake/internal/game"
	"github.com/C0d3-5t3w/go-snake/internal/gui"
	"github.com/C0d3-5t3w/go-snake/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize storage
	store, err := storage.NewStorage()
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Create game instance
	gameInstance := game.NewGame(cfg)

	// Create Ebiten GUI wrapper
	ebitenGUI, err := gui.NewEbitenGUI(gameInstance, cfg, store)
	if err != nil {
		log.Fatalf("Failed to initialize Ebiten GUI: %v", err)
	}

	// Run the game using Ebiten's RunGame function
	if err := ebitenGUI.Run(); err != nil {
		log.Fatalf("Ebiten RunGame error: %v", err)
	}
}

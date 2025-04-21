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

	// Create GUI
	guiInstance, err := gui.NewGUI(gameInstance, cfg, store)
	if err != nil {
		log.Fatalf("Failed to initialize GUI: %v", err)
	}

	// Run the game
	guiInstance.Run()
}

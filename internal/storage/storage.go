package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// HighScore represents a player's high score
type HighScore struct {
	Player string    `json:"player"`
	Score  int       `json:"score"`
	Date   time.Time `json:"date"`
}

// Settings represents user settings
type Settings struct {
	MusicVolume float64 `json:"music_volume"`
	SfxVolume   float64 `json:"sfx_volume"`
	Difficulty  string  `json:"difficulty"`
}

// GameData represents all persistent game data
type GameData struct {
	HighScores []HighScore `json:"high_scores"`
	Settings   Settings    `json:"settings"`
}

// Storage handles game data persistence
type Storage struct {
	filePath string
	data     GameData
}

// NewStorage creates a new storage instance
func NewStorage() (*Storage, error) {
	storagePath := findStoragePath()
	storage := &Storage{filePath: storagePath}

	// Try to load existing data, or create default if file doesn't exist
	if err := storage.Load(); err != nil {
		// Create default data
		storage.data = GameData{
			HighScores: []HighScore{},
			Settings: Settings{
				MusicVolume: 0.7,
				SfxVolume:   0.8,
				Difficulty:  "medium",
			},
		}
		// Save default data
		if err := storage.Save(); err != nil {
			return nil, err
		}
	}

	return storage, nil
}

// Load reads data from storage file
func (s *Storage) Load() error {
	data, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.data)
}

// Save writes data to storage file
func (s *Storage) Save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.filePath, data, 0644)
}

// AddHighScore adds a new high score and maintains order
func (s *Storage) AddHighScore(player string, score int) {
	newScore := HighScore{
		Player: player,
		Score:  score,
		Date:   time.Now(),
	}

	s.data.HighScores = append(s.data.HighScores, newScore)

	// Sort high scores in descending order
	sort.Slice(s.data.HighScores, func(i, j int) bool {
		return s.data.HighScores[i].Score > s.data.HighScores[j].Score
	})

	// Keep only top 10 scores
	if len(s.data.HighScores) > 10 {
		s.data.HighScores = s.data.HighScores[:10]
	}
}

// GetHighScores returns all high scores
func (s *Storage) GetHighScores() []HighScore {
	return s.data.HighScores
}

// GetSettings returns the current settings
func (s *Storage) GetSettings() Settings {
	return s.data.Settings
}

// UpdateSettings updates the game settings
func (s *Storage) UpdateSettings(settings Settings) {
	s.data.Settings = settings
}

// findStoragePath locates the storage.json file
func findStoragePath() string {
	// Try different common locations
	locations := []string{
		filepath.Join("pkg", "storage", "storage.json"),
		filepath.Join("..", "pkg", "storage", "storage.json"),
		filepath.Join("..", "..", "pkg", "storage", "storage.json"),
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc
		}
	}

	// Default location if not found
	return filepath.Join("pkg", "storage", "storage.json")
}

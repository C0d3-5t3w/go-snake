package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the game configuration
type Config struct {
	Game struct {
		GridSize       int     `yaml:"grid_size"`
		InitialSpeed   float64 `yaml:"initial_speed"`
		SpeedIncrement float64 `yaml:"speed_increment"`
		MaxSpeed       float64 `yaml:"max_speed"`
		InitialLength  int     `yaml:"initial_length"`
	} `yaml:"game"`

	Graphics struct {
		WindowWidth  int  `yaml:"window_width"`
		WindowHeight int  `yaml:"window_height"`
		Fullscreen   bool `yaml:"fullscreen"`
		Vsync        bool `yaml:"vsync"`
	} `yaml:"graphics"`

	Controls struct {
		Forward  string `yaml:"forward"`
		Backward string `yaml:"backward"`
		Left     string `yaml:"left"`
		Right    string `yaml:"right"`
		Pause    string `yaml:"pause"`
		Quit     string `yaml:"quit"`
	} `yaml:"controls"`

	Colors struct {
		SnakeHead  [3]float32 `yaml:"snake_head"`
		SnakeBody  [3]float32 `yaml:"snake_body"`
		Food       [3]float32 `yaml:"food"`
		Grid       [3]float32 `yaml:"grid"`
		Background [3]float32 `yaml:"background"`
	} `yaml:"colors"`
}

// LoadConfig loads configuration from the config file
func LoadConfig() (*Config, error) {
	// Find the config file
	configPath := findConfigPath()

	// Read the config file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse the config
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// findConfigPath locates the config.yaml file
func findConfigPath() string {
	// Try different common locations
	locations := []string{
		filepath.Join("pkg", "config", "config.yaml"),
		filepath.Join("..", "pkg", "config", "config.yaml"),
		filepath.Join("..", "..", "pkg", "config", "config.yaml"),
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc
		}
	}

	// Default location if not found
	log.Println("Config file not found in standard locations, using default path")
	return filepath.Join("pkg", "config", "config.yaml")
}

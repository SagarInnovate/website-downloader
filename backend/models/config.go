package models

import (
	"encoding/json"
	"os"
)

// Config represents the application configuration
type Config struct {
	MaxDepth          int      `json:"maxDepth"`
	MaxPages          int      `json:"maxPages"`
	WorkerCount       int      `json:"workerCount"`
	RequestTimeout    int      `json:"requestTimeout"`
	RateLimit         int      `json:"rateLimit"`
	AllowedExtensions []string `json:"allowedExtensions"`
	TempDir           string   `json:"tempDir"`
	OutputDir         string   `json:"outputDir"`
	UserAgent         string   `json:"userAgent"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		MaxDepth:       5,
		MaxPages:       1000,
		WorkerCount:    10,
		RequestTimeout: 30,
		RateLimit:      5,
		AllowedExtensions: []string{
			".html", ".htm", ".css", ".js", ".jsx", ".ts", ".tsx",
			".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp", ".ico",
			".woff", ".woff2", ".ttf", ".eot",
			".pdf", ".json", ".xml",
		},
		TempDir:   "./temp",
		OutputDir: "./downloads",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}
}

// LoadConfig loads configuration from a file
func LoadConfig(filepath string) (*Config, error) {
	// Try to load from file
	data, err := os.ReadFile(filepath)
	if err != nil {
		// If file doesn't exist, return default config
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(filepath string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

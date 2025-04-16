package internal

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".relaise", "config.yaml"), nil
}

func LoadConfig() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func DefaultConfig() *Config {
	return &Config{
		APIKey:          "",
		Provider:        "mistral",
		Model:           "mistral-small-latest",
		Mood:            "professional",
		ReleaseType:     "minor",
		BulletStyle:     "-",
		IncludeSections: false,
		Language:        "en",
		Emojis:          false,
		Copy:            false,
		Temperature:     0.3,
		MaxTokens:       1000,
	}
}

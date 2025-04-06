package services

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIKey          string `yaml:"api_key"`
	Provider        string `yaml:"provider"`
	Model           string `yaml:"model"`
	Mood            string `yaml:"mood"`
	ReleaseType     string `yaml:"release_type"`     // minor, major, patch
	BulletStyle     string `yaml:"bullet_style"`     // "*", "-", or numbered
	IncludeSections bool   `yaml:"include_sections"` // Features, Fixes, etc.
	Language        string `yaml:"language"`         // e.g., "en", "fr", "es"
}

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

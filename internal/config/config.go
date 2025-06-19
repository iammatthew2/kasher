package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type TaskConfig struct {
	Command     string `toml:"command"`
	Expiration  string `toml:"expiration"`
	Notes       string `toml:"notes,omitempty"`
	LastFetched string `toml:"lastFetched,omitempty"`
}

type KasherConfig map[string]TaskConfig

// getConfigPath returns the path to the kasher config file.
// On macOS, this will be "$HOME/Library/Application Support/kasher/config.toml".
// It creates the kasher config directory if it does not exist.
func getConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	kasherDir := filepath.Join(dir, "kasher")
	err = os.MkdirAll(kasherDir, 0o755)
	if err != nil {
		return "", err
	}
	return filepath.Join(kasherDir, "config.toml"), nil
}

// LoadConfig loads the kasher configuration from disk.
// If the config file does not exist, it returns an empty KasherConfig.
func LoadConfig() (KasherConfig, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return make(KasherConfig), nil // return empty config
	} else if err != nil {
		return nil, err
	}
	var cfg KasherConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// SaveConfig writes the provided KasherConfig to disk in TOML format.
// It overwrites the existing config file or creates a new one if it does not exist.
func SaveConfig(cfg KasherConfig) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// AddTask adds a new task to the config. Returns an error if the task already exists.
func (cfg KasherConfig) AddTask(name string, task TaskConfig) error {
	if _, exists := cfg[name]; exists {
		return errors.New("task already exists")
	}
	cfg[name] = task
	return nil
}

// UpdateTask updates an existing task in the config. Returns an error if the task does not exist.
func (cfg KasherConfig) UpdateTask(name string, task TaskConfig) error {
	if _, exists := cfg[name]; !exists {
		return errors.New("task does not exist")
	}
	cfg[name] = task
	return nil
}

// DeleteTask removes a task from the config. Returns an error if the task does not exist.
func (cfg KasherConfig) DeleteTask(name string) error {
	if _, exists := cfg[name]; !exists {
		return errors.New("task does not exist")
	}
	delete(cfg, name)
	return nil
}

// ClearConfig deletes the kasher config file from disk.
// Returns nil if successful, or an error if the file could not be deleted.
func ClearConfig() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	// Ignore error if file does not exist
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// GetConfigPath returns the path to the kasher config file.
func GetConfigPath() (string, error) {
	return getConfigPath()
}

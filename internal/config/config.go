package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type TaskConfig struct {
	Command    string `toml:"command"`
	Expiration string `toml:"expiration"`
	Notes      string `toml:"notes,omitempty"`
}

type KasherConfig map[string]TaskConfig

// getConfigPath returns ~/.config/kasher/config.toml
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

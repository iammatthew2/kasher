package config

import (
	"os"
	"path/filepath"
)

// GetCacheFilePath returns the path to the cache file for a given task.
func GetCacheFilePath(taskName string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(dir, "kasher")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", err
	}
	cachePath := filepath.Join(cacheDir, taskName+".cache")
	return cachePath, nil
}

// WriteCache saves the output to the cache file for the given task.
func WriteCache(taskName, output string) error {
	path, err := GetCacheFilePath(taskName)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(output), 0o644)
}

// ReadCache reads the cached output for the given task.
func ReadCache(taskName string) (string, error) {
	path, err := GetCacheFilePath(taskName)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

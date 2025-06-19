package config

import (
	"log"
	"os"
	"path/filepath"
)

// GetCacheFilePath returns the path to the cache file for a given task.
func GetCacheFilePath(taskName string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		log.Printf("[kasher] Error getting user cache dir: %v", err)
		return "", err
	}
	cacheDir := filepath.Join(dir, "kasher")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		log.Printf("[kasher] Error creating cache dir %s: %v", cacheDir, err)
		return "", err
	}
	cachePath := filepath.Join(cacheDir, taskName+".cache")
	log.Printf("[kasher] Cache file path for task '%s': %s", taskName, cachePath)
	return cachePath, nil
}

// WriteCache saves the output to the cache file for the given task.
func WriteCache(taskName, output string) error {
	path, err := GetCacheFilePath(taskName)
	if err != nil {
		log.Printf("[kasher] Error getting cache file path for write: %v", err)
		return err
	}
	err = os.WriteFile(path, []byte(output), 0o644)
	if err != nil {
		log.Printf("[kasher] Error writing cache file %s: %v", path, err)
	} else {
		log.Printf("[kasher] Successfully wrote cache file: %s", path)
	}
	return err
}

// ReadCache reads the cached output for the given task.
func ReadCache(taskName string) (string, error) {
	path, err := GetCacheFilePath(taskName)
	if err != nil {
		log.Printf("[kasher] Error getting cache file path for read: %v", err)
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("[kasher] Error reading cache file %s: %v", path, err)
		return "", err
	}
	log.Printf("[kasher] Successfully read cache file: %s", path)
	return string(data), nil
}

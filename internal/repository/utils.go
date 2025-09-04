package repository

import (
	"fmt"
	"os"
	"path/filepath"
)

// findGitDir recursively searches for a .git directory upwards from the current directory
func findGitDir(path string) (string, error) {
	for {
		// Check if the .git directory exists in the current path
		gitPath := filepath.Join(path, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return gitPath, nil
		}

		// Check if we've reached the root directory
		parent := filepath.Dir(path)
		if parent == path {
			break
		}

		// Move up one directory
		path = parent
	}

	return "", fmt.Errorf(".git directory not found")
}

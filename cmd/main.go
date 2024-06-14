package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/plumbing/object"
	"github.com/go-git/go-git/v5"
)

func main() {
	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		slog.Error("could not get cwd", "error", err)
		return
	}

	// Search for the .git directory
	gitDir, err := findGitDir(currentDir)
	if err != nil {
		slog.Error("could not find .git directory", "error", err)
		return
	}

	// open repository
	r, err := git.PlainOpen(gitDir)
	if err != nil {
		slog.Error("could not open repository", "error", err)
		return
	}

	// get the HEAD reference
	hd, err := r.Head()
	if err != nil {
		slog.Error("could not get HEAD reference", "error", err)
		return
	}

	// get the commit log iterator
	citr, err := r.Log(&git.LogOptions{From: hd.Hash()})
	if err != nil {
		slog.Error("could not get commit log", "error", err)
		return
	}

	err = citr.ForEach(func(c *object.Commit) error {
		slog.Info("commit", "hash", c.Hash, "message", c.Message)
		return nil
	})
}

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

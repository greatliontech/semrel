package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/greatliontech/semrel"
	"github.com/greatliontech/semrel/internal/cmd"
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

	// get config
	cfg, err := semrel.ParseConfigFile(filepath.Join(filepath.Dir(gitDir), ".semrel"))
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Error("could not parse config file", "error", err)
			return
		}
		slog.Debug("config file not found, using default config")
		cfg = semrel.DefaultConfig
	}

	// create cli
	cli, err := cmd.New(r, cfg)
	if err != nil {
		slog.Error("could not create CLI", "error", err)
		return
	}

	// run the CLI
	if err := cli.Execute(); err != nil {
		slog.Error("could not execute CLI", "error", err)
		return
	}
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

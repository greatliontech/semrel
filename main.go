package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/greatliontech/semrel/internal/cmd"
	"github.com/greatliontech/semrel/internal/repository"
	"github.com/greatliontech/semrel/pkg/semrel"
)

var version = "0.0.0-dev"

func main() {
	// open repository
	r, err := repository.Open()
	if err != nil {
		slog.Error("could not open repository", "error", err)
		os.Exit(1)
	}

	var cfg *semrel.Config

	// get config
	cfgFile, err := semrel.ConfigFileFromPath(filepath.Join(r.Root(), ".semrel.yaml"))
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Error("could not parse config file", "error", err)
			os.Exit(1)
		}
		slog.Debug("config file not found, using default config")
		cfg = semrel.DefaultConfig
	} else {
		cfg, err = semrel.NewConfigFromConfigFile(cfgFile)
		if err != nil {
			slog.Error("could not create config from config file", "error", err)
			os.Exit(1)
		}
	}

	// create cli
	cli, err := cmd.New(r, cfg, version)
	if err != nil {
		slog.Error("could not create CLI", "error", err)
		os.Exit(1)
	}

	// run the CLI
	code := cli.Execute()
	os.Exit(code)
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/greatliontech/semrel/internal/repository"
	"github.com/spf13/cobra"
)

const configFileName = ".semrel.yaml"

const schemaComment = "# yaml-language-server: $schema=https://raw.githubusercontent.com/greatliontech/semrel/refs/heads/main/config.schema.json"

const fullConfig = `# yaml-language-server: $schema=https://raw.githubusercontent.com/greatliontech/semrel/refs/heads/main/config.schema.json
version: v0

# Commit types that trigger a patch version bump
patchTypes:
  - fix

# Commit types that trigger a minor version bump
minorTypes:
  - feat

# Commit types that trigger a major version bump
# majorTypes:
#   - breaking

# Default bump type when no commit types match (none, patch, minor, major)
# defaultBump: none

# Initial version when no tags exist (default: 1.0.0)
# initialVersion: 1.0.0

# Development mode: starts at 0.1.0, treats breaking changes as minor bumps
# development: false

# Version prefix (e.g., "v" for v1.0.0)
# prefix: ""

# Automatically create git tags
# createTag: false

# Automatically push tags to remote (requires createTag)
# pushTag: false

# Release platform (github, gitlab)
# platform: ""

# Regex rules for transforming commit messages in release notes
# matchRules:
#   - match: "TICKET-(\\d+)"
#     replace: "[#$1](https://tracker.example.com/TICKET-$1)"

# Exclude certain commit types or scopes from release notes
# filters:
#   types:
#     - docs
#     - test
#   scopes:
#     - internal
`

type initCommand struct {
	cmd         *cobra.Command
	repo        *repository.Repo
	force       bool
	development bool
	prefix      string
	full        bool
}

func newInitCommand(repo *repository.Repo) *initCommand {
	c := &initCommand{
		repo: repo,
	}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new .semrel.yaml configuration file",
		RunE:  c.runE,
	}
	cmd.Flags().BoolVarP(&c.force, "force", "f", false, "overwrite existing config file")
	cmd.Flags().BoolVarP(&c.development, "development", "d", false, "enable development mode (0.x.x versions)")
	cmd.Flags().StringVarP(&c.prefix, "prefix", "p", "", "version prefix (e.g., 'v')")
	cmd.Flags().BoolVar(&c.full, "full", false, "generate full commented config with all options")
	c.cmd = cmd
	return c
}

func (c *initCommand) runE(cmd *cobra.Command, args []string) error {
	configPath := filepath.Join(c.repo.Root(), configFileName)

	// Check if file already exists
	if _, err := os.Stat(configPath); err == nil {
		if !c.force {
			return fmt.Errorf("config file already exists: %s (use --force to overwrite)", configPath)
		}
	}

	config := buildInitConfig(c.development, c.prefix, c.full)

	// Write the config file
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Created %s\n", configPath)
	return nil
}

func buildInitConfig(development bool, prefix string, full bool) string {
	if full {
		return fullConfig
	}

	var b strings.Builder
	b.WriteString(schemaComment)
	b.WriteString("\nversion: v0\n")

	if development {
		b.WriteString("development: true\n")
	}

	if prefix != "" {
		b.WriteString(fmt.Sprintf("prefix: %q\n", prefix))
	}

	return b.String()
}

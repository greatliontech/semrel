package cmd

import (
	"context"
	"log"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/greatliontech/semrel/internal/repository"
	"github.com/greatliontech/semrel/pkg/semrel"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

type mcpCommand struct {
	cmd     *cobra.Command
	repo    *repository.Repo
	cfg     *semrel.Config
	version string
}

func newMCPCommand(repo *repository.Repo, cfg *semrel.Config, version string) *mcpCommand {
	c := &mcpCommand{
		repo:    repo,
		cfg:     cfg,
		version: version,
	}
	c.cmd = &cobra.Command{
		Use:   "mcp",
		Short: "Run as MCP server",
		Long:  "Run semrel as a Model Context Protocol server for LLM integration",
		RunE:  c.runE,
	}
	return c
}

func (c *mcpCommand) runE(cmd *cobra.Command, args []string) error {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "semrel",
		Version: c.version,
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "config",
		Description: "Get the semrel configuration including development mode and commit type mappings",
	}, c.handleConfig)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "current_version",
		Description: "Get the current semantic version from git tags",
	}, c.handleCurrentVersion)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "next_version",
		Description: "Calculate the next semantic version based on conventional commits since the last release",
	}, c.handleNextVersion)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
	return nil
}

// ConfigArgs are the arguments for the config tool.
type ConfigArgs struct{}

// ConfigResult is the result of the config tool.
type ConfigResult struct {
	Development bool     `json:"development"`
	PatchTypes  []string `json:"patch_types"`
	MinorTypes  []string `json:"minor_types"`
	MajorTypes  []string `json:"major_types"`
	Prefix      string   `json:"prefix"`
}

func (c *mcpCommand) handleConfig(ctx context.Context, req *mcp.CallToolRequest, args ConfigArgs) (*mcp.CallToolResult, *ConfigResult, error) {
	result := &ConfigResult{
		Development: c.cfg.IsDevelopment(),
		PatchTypes:  c.cfg.PatchTypes(),
		MinorTypes:  c.cfg.MinorTypes(),
		MajorTypes:  c.cfg.MajorTypes(),
		Prefix:      c.cfg.Prefix(),
	}
	return &mcp.CallToolResult{}, result, nil
}

// CurrentVersionArgs are the arguments for the current_version tool.
type CurrentVersionArgs struct{}

// CurrentVersionResult is the result of the current_version tool.
type CurrentVersionResult struct {
	Version string `json:"version"`
}

func (c *mcpCommand) handleCurrentVersion(ctx context.Context, req *mcp.CallToolRequest, args CurrentVersionArgs) (*mcp.CallToolResult, *CurrentVersionResult, error) {
	ver, _, err := c.repo.CurrentVersion(false, c.cfg.StableOnly())
	if err != nil {
		return nil, nil, err
	}

	result := &CurrentVersionResult{
		Version: ver.String(),
	}

	return &mcp.CallToolResult{}, result, nil
}

// NextVersionArgs are the arguments for the next_version tool.
type NextVersionArgs struct{}

// NextVersionResult is the result of the next_version tool.
type NextVersionResult struct {
	Current  string `json:"current"`
	Next     string `json:"next"`
	BumpType string `json:"bump_type"`
}

func (c *mcpCommand) handleNextVersion(ctx context.Context, req *mcp.CallToolRequest, args NextVersionArgs) (*mcp.CallToolResult, *NextVersionResult, error) {
	currentVer, ref, err := c.repo.CurrentVersion(false, c.cfg.StableOnly())
	if err != nil {
		return nil, nil, err
	}

	head, err := c.repo.Head()
	if err != nil {
		return nil, nil, err
	}

	var toHash plumbing.Hash
	if ref != nil {
		toHash = ref.Hash()
	}

	commits, err := c.repo.Commits(head, toHash)
	if err != nil {
		return nil, nil, err
	}

	// Calculate bump kind
	bumpKind := semrel.BumpNone
	for _, commit := range commits {
		b := commit.BumpKind(c.cfg)
		if b.IsGreater(bumpKind) {
			bumpKind = b
		}
		if bumpKind == semrel.BumpMajor {
			break
		}
	}

	nextVer := semrel.NextVersion(currentVer, commits, c.cfg)

	result := &NextVersionResult{
		Current:  currentVer.String(),
		Next:     nextVer.String(),
		BumpType: bumpKind.String(),
	}

	return &mcp.CallToolResult{}, result, nil
}

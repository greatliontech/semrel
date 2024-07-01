package cmd

import (
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

type currentCommand struct {
	cmd  *cobra.Command
	repo *git.Repository
}

func newCurrentCommand(repo *git.Repository) *currentCommand {
	c := &currentCommand{
		repo: repo,
	}
	cmd := &cobra.Command{
		Use:  "current",
		RunE: c.runE,
	}
	c.cmd = cmd
	return c
}

func (c *currentCommand) runE(cmd *cobra.Command, args []string) error {
	// get the HEAD reference
	hd, err := c.repo.Head()
	if err != nil {
		return err
	}
	cmd.Println(hd.Hash().String())
	return nil
}

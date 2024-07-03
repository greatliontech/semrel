package cmd

import (
	"fmt"

	"github.com/greatliontech/semrel/internal/repository"
	"github.com/spf13/cobra"
)

type currentCommand struct {
	cmd  *cobra.Command
	repo *repository.Repo
}

func newCurrentCommand(repo *repository.Repo) *currentCommand {
	c := &currentCommand{
		repo: repo,
	}
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Print the current release version",
		RunE:  c.runE,
	}
	c.cmd = cmd
	return c
}

func (c *currentCommand) runE(cmd *cobra.Command, args []string) error {
	cv, _, err := c.repo.CurrentVersion()
	if err != nil {
		return err
	}
	fmt.Println(cv.String())
	return nil
}

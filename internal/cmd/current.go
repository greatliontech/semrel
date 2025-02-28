package cmd

import (
	"fmt"

	"github.com/greatliontech/semrel/internal/repository"
	"github.com/greatliontech/semrel/pkg/semrel"
	"github.com/spf13/cobra"
)

type currentCommand struct {
	cmd               *cobra.Command
	repo              *repository.Repo
	cfg               *semrel.Config
	currentBranchOnly bool
}

func newCurrentCommand(repo *repository.Repo, cfg *semrel.Config) *currentCommand {
	c := &currentCommand{
		repo: repo,
		cfg:  cfg,
	}
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Print the current release version",
		RunE:  c.runE,
	}
	cmd.Flags().BoolVarP(&c.currentBranchOnly, "current-branch-only", "", false, "only tags from the current branch")
	c.cmd = cmd
	return c
}

func (c *currentCommand) runE(cmd *cobra.Command, args []string) error {
	cv, _, err := c.repo.CurrentVersion(c.currentBranchOnly)
	if err != nil {
		return err
	}
	currentTag := fmt.Sprintf("%s%s", c.cfg.Prefix(), cv.String())
	fmt.Println(currentTag)
	return nil
}

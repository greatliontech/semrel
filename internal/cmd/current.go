package cmd

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/greatliontech/semrel/internal/repository"
	"github.com/spf13/cobra"
)

type currentCommand struct {
	cmd               *cobra.Command
	repo              *repository.Repo
	currentBranchOnly bool
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
	cmd.Flags().BoolVarP(&c.currentBranchOnly, "current-branch-only", "", false, "only tags from the current branch")
	c.cmd = cmd
	return c
}

var emptyVersion = semver.New(0, 0, 0, "", "")

func (c *currentCommand) runE(cmd *cobra.Command, args []string) error {
	cv, _, err := c.repo.CurrentVersion(c.currentBranchOnly)
	if err != nil {
		if err == repository.ErrNoTags {
			fmt.Println(emptyVersion.String())
			return nil
		}
		return err
	}
	fmt.Println(cv.String())
	return nil
}

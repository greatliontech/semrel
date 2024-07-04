package cmd

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/greatliontech/semrel/internal/repository"
	"github.com/spf13/cobra"
)

type compareCommand struct {
	cmd               *cobra.Command
	repo              *repository.Repo
	le                []string
	ge                []string
	lt                []string
	gt                []string
	currentBranchOnly bool
}

func newCompareCommand(repo *repository.Repo) *compareCommand {
	c := &compareCommand{
		repo: repo,
	}
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare the current or supplied version with the given versions",
		RunE:  c.runE,
		Args:  cobra.MaximumNArgs(1),
	}
	cmd.Flags().StringSliceVarP(&c.le, "le", "", nil, "less than or equal to")
	cmd.Flags().StringSliceVarP(&c.ge, "ge", "", nil, "greater than or equal to")
	cmd.Flags().StringSliceVarP(&c.lt, "lt", "", nil, "less than")
	cmd.Flags().StringSliceVarP(&c.gt, "gt", "", nil, "greater than")
	cmd.Flags().BoolVarP(&c.currentBranchOnly, "current-branch-only", "", false, "only compare the current branch")
	c.cmd = cmd
	return c
}

var errCompareFailed = errors.New("compare failed")

func (c *compareCommand) runE(cmd *cobra.Command, args []string) error {
	var current *semver.Version
	var err error
	if len(args) == 0 {
		current, _, err = c.repo.CurrentVersion(c.currentBranchOnly)
	} else {
		current, err = semver.NewVersion(args[0])
	}
	if err != nil {
		return err
	}

	for _, v := range c.le {
		other, err := semver.NewVersion(v)
		if err != nil {
			return err
		}
		if !current.LessThan(other) && !current.Equal(other) {
			return errCompareFailed
		}
	}

	for _, v := range c.ge {
		other, err := semver.NewVersion(v)
		if err != nil {
			return err
		}
		if !current.GreaterThan(other) && !current.Equal(other) {
			return errCompareFailed
		}
	}

	for _, v := range c.lt {
		other, err := semver.NewVersion(v)
		if err != nil {
			return err
		}
		if !current.LessThan(other) {
			return errCompareFailed
		}
	}

	for _, v := range c.gt {
		other, err := semver.NewVersion(v)
		if err != nil {
			return err
		}
		if !current.GreaterThan(other) {
			return errCompareFailed
		}
	}

	return nil
}

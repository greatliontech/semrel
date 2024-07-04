package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/greatliontech/semrel/internal/repository"
	"github.com/greatliontech/semrel/pkg/semrel"
	"github.com/spf13/cobra"
)

type rootCommand struct {
	cmd       *cobra.Command
	repo      *repository.Repo
	cfg       *semrel.Config
	createTag bool
	pushTag   bool
}

func New(r *git.Repository, cfg *semrel.Config, ver string) (*rootCommand, error) {
	rp := repository.New(r)
	c := &rootCommand{
		repo: rp,
		cfg:  cfg,
	}
	cmd := &cobra.Command{
		Use:           "semrel",
		RunE:          c.runE,
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       ver,
	}
	cmd.Flags().BoolVarP(&c.createTag, "create-tag", "", false, "create the tag")
	cmd.Flags().BoolVarP(&c.pushTag, "push-tag", "", false, "push the tag")
	cmd.AddCommand(
		newCurrentCommand(rp).cmd,
		newCompareCommand(rp).cmd,
		newValidateCommand().cmd,
	)
	c.cmd = cmd
	return c, nil
}

func (r *rootCommand) Execute() {
	err := r.cmd.Execute()
	if err != nil {
		if err != errCompareFailed {
			slog.Error("command failed", "error", err)
		}
		os.Exit(1)
	}
	os.Exit(0)
}

func (r *rootCommand) runE(cmd *cobra.Command, args []string) error {
	// get latest tag version
	ver, ref, err := r.repo.CurrentVersion()
	if err != nil {
		if err == repository.ErrNoTags {
			fmt.Println(r.cfg.InitialVersion().String())
			return nil
		}
		return err
	}

	commits, err := r.repo.Commits(plumbing.ZeroHash, ref.Hash())
	if err != nil {
		return err
	}

	next := semrel.NextVersion(ver, commits, r.cfg)

	nextTag := fmt.Sprintf("%s%s", r.cfg.Prefix(), next.String())

	if r.createTag {
		head, err := r.repo.Head()
		if err != nil {
			return err
		}
		err = r.repo.CreateTag(nextTag, head, r.pushTag)
		if err != nil {
			return err
		}
	}

	fmt.Println(nextTag)

	return nil
}

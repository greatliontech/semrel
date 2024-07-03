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
	cmd  *cobra.Command
	repo *repository.Repo
	cfg  *semrel.Config
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

	fmt.Println(next.String())

	return nil
}

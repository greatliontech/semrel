package cmd

import (
	"log/slog"

	"github.com/go-git/go-git/plumbing/object"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/thegrumpylion/semrel"
)

type rootCommand struct {
	cmd  *cobra.Command
	repo *git.Repository
}

func New(repo *git.Repository) (*rootCommand, error) {
	r := &rootCommand{}
	cmd := &cobra.Command{
		Use:  "semrel",
		RunE: r.runE,
	}
	r.cmd = cmd
	return r, nil
}

func (r *rootCommand) Execute() error {
	return r.cmd.Execute()
}

func (r *rootCommand) runE(cmd *cobra.Command, args []string) error {
	// get the HEAD reference
	hd, err := r.repo.Head()
	if err != nil {
		slog.Error("could not get HEAD reference", "error", err)
		return err
	}

	// get the commit log iterator
	citr, err := r.Log(&git.LogOptions{From: hd.Hash()})
	if err != nil {
		slog.Error("could not get commit log", "error", err)
		return err
	}

	err = citr.ForEach(func(c *object.Commit) error {
		semrel.ParseCommitMessage(c.Message)
		return nil
	})

	return nil
}

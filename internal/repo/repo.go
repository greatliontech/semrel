package repo

import (
	"errors"
	"log/slog"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/greatliontech/semrel"
)

type Repo struct {
	repo *git.Repository
}

func New(repo *git.Repository) *Repo {
	return &Repo{repo: repo}
}

func (r *Repo) Commits(from, to plumbing.Hash) ([]*semrel.Commit, error) {
	// get the commit log iterator
	citr, err := r.repo.Log(&git.LogOptions{From: from})
	if err != nil {
		slog.Error("could not get commit log", "error", err)
		return nil, err
	}

	errBreak := errors.New("break")
	commits := []*semrel.Commit{}
	err = citr.ForEach(func(c *object.Commit) error {
		if c.Hash == to {
			return errBreak
		}
		cmt, err := semrel.ParseCommitMessage(c.Message)
		if err != nil {
			if err == semrel.ErrNotConventionalCommit {
				return nil
			}
			return err
		}
		commits = append(commits, cmt)
		return nil
	})
	if err != nil && err != errBreak {
		slog.Error("could not iterate over commits", "error", err)
		return nil, err
	}

	return commits, nil
}

package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/greatliontech/semrel/pkg/semrel"
)

type Repo struct {
	repo *git.Repository
}

func New(repo *git.Repository) *Repo {
	return &Repo{repo: repo}
}

func (r *Repo) Head() (plumbing.Hash, error) {
	ref, err := r.repo.Head()
	if err != nil {
		return plumbing.ZeroHash, err
	}
	return ref.Hash(), nil
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

type versionReference struct {
	ver *semver.Version
	ref *plumbing.Reference
}

var ErrNoTags = errors.New("no tags found")

func (r *Repo) CurrentVersion() (*semver.Version, *plumbing.Reference, error) {
	// get the tag iterator
	titr, err := r.repo.Tags()
	if err != nil {
		return nil, nil, err
	}
	versions := []versionReference{}
	err = titr.ForEach(func(ref *plumbing.Reference) error {
		ver, err := semver.NewVersion(ref.Name().Short())
		if err == nil {
			versions = append(versions, versionReference{ver: ver, ref: ref})
		}
		return nil
	})
	sort.Slice(versions, func(i, j int) bool {
		v1 := versions[i].ver
		v2 := versions[j].ver
		return v2.LessThan(v1)
	})
	if len(versions) == 0 {
		return nil, nil, ErrNoTags
	}
	return versions[0].ver, versions[0].ref, nil
}

func (r *Repo) CreateTag(tag string, commit plumbing.Hash, push bool) error {
	_, err := r.repo.CreateTag(tag, commit, nil)
	if err != nil {
		return err
	}

	if push {

		refSpecStr := fmt.Sprintf("refs/tags/%s:refs/tags/%s", tag, tag)
		refSpec := config.RefSpec(refSpecStr)

		return r.repo.Push(&git.PushOptions{
			RemoteName: "origin",
			RefSpecs: []config.RefSpec{
				refSpec,
			},
		})
	}

	return err
}

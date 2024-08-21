package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"

	"github.com/Masterminds/semver/v3"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
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
	citr, err := r.repo.Log(&git.LogOptions{
		From:  from,
		Order: git.LogOrderDFSPost,
	})
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

var emptyVersion = semver.New(0, 0, 0, "", "")

type versionReference struct {
	ver *semver.Version
	ref *plumbing.Reference
}

func (r *Repo) CurrentVersion(currentBranchOnly bool) (*semver.Version, *plumbing.Reference, error) {
	currentBranchRefs := mapset.NewSet[plumbing.Hash]()

	if currentBranchOnly {
		head, err := r.repo.Head()
		if err != nil {
			return nil, nil, err
		}
		litr, err := r.repo.Log(&git.LogOptions{From: head.Hash()})
		if err != nil {
			return nil, nil, err
		}
		err = litr.ForEach(func(c *object.Commit) error {
			currentBranchRefs.Add(c.Hash)
			return nil
		})
		if err != nil {
			return nil, nil, err
		}
	}

	// get the tag iterator
	titr, err := r.repo.Tags()
	if err != nil {
		return nil, nil, err
	}
	versions := []versionReference{}

	err = titr.ForEach(func(ref *plumbing.Reference) error {
		if currentBranchOnly {
			if !currentBranchRefs.Contains(ref.Hash()) {
				return nil
			}
		}
		ver, err := semver.NewVersion(ref.Name().Short())
		if err == nil {
			versions = append(versions, versionReference{ver: ver, ref: ref})
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	// sort the versions in descending order
	sort.Slice(versions, func(i, j int) bool {
		v1 := versions[i].ver
		v2 := versions[j].ver
		return v2.LessThan(v1)
	})
	if len(versions) == 0 {
		return emptyVersion, nil, nil
	}
	return versions[0].ver, versions[0].ref, nil
}

func (r *Repo) CreateTag(tag string, commit plumbing.Hash, push bool, auth transport.AuthMethod) error {
	_, err := r.repo.CreateTag(tag, commit, nil)
	if err != nil {
		return err
	}

	if push {
		refSpecStr := fmt.Sprintf("refs/tags/%s:refs/tags/%s", tag, tag)
		refSpec := config.RefSpec(refSpecStr)

		return r.repo.Push(&git.PushOptions{
			RefSpecs: []config.RefSpec{
				refSpec,
			},
			Auth: auth,
		})
	}

	return err
}

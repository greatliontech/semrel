package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/greatliontech/semrel"
	"github.com/spf13/cobra"
)

type rootCommand struct {
	cmd  *cobra.Command
	repo *git.Repository
	cfg  *semrel.Config
}

func New(repo *git.Repository, cfg *semrel.Config) (*rootCommand, error) {
	r := &rootCommand{
		repo: repo,
		cfg:  cfg,
	}
	cmd := &cobra.Command{
		Use:  "semrel",
		RunE: r.runE,
	}
	cmd.AddCommand(
		newTagsCommand(repo).cmd,
		newCurrentCommand(repo).cmd,
	)
	r.cmd = cmd
	return r, nil
}

func (r *rootCommand) Execute() error {
	return r.cmd.Execute()
}

func (r *rootCommand) runE(cmd *cobra.Command, args []string) error {
	// get latest tag version
	ver, ref, err := getLatestTagVersion(r.repo)
	if err != nil {
		if err == errNoTags {
			fmt.Println("0.0.1")
			return nil
		}
		return err
	}

	// get the HEAD reference
	hd, err := r.repo.Head()
	if err != nil {
		slog.Error("could not get HEAD reference", "error", err)
		return err
	}

	// get the commit log iterator
	citr, err := r.repo.Log(&git.LogOptions{From: hd.Hash()})
	if err != nil {
		slog.Error("could not get commit log", "error", err)
		return err
	}

	errBreak := errors.New("break")
	commits := []*semrel.Commit{}
	err = citr.ForEach(func(c *object.Commit) error {
		if c.Hash == ref.Hash() {
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
		return err
	}

	if len(commits) == 0 {
		fmt.Println(ver.String())
	}

	next := semrel.NextVersion(ver, commits, semrel.DefaultConfig)

	fmt.Println(next.String())

	return nil
}

type versionReference struct {
	ver *semver.Version
	ref *plumbing.Reference
}

var errNoTags = errors.New("no tags found")

func getLatestTagVersion(repo *git.Repository) (*semver.Version, *plumbing.Reference, error) {
	// get the tag iterator
	titr, err := repo.Tags()
	if err != nil {
		return nil, nil, err
	}
	versions := []versionReference{}
	err = titr.ForEach(func(ref *plumbing.Reference) error {
		ver, err := semver.NewVersion(ref.Name().Short())
		if err == nil && ver.Prerelease() == "" {
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
		return nil, nil, errNoTags
	}
	return versions[0].ver, versions[0].ref, nil
}

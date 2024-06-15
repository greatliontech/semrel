package cmd

import (
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

type tagsCommand struct {
	cmd  *cobra.Command
	repo *git.Repository
}

func newTagsCommand(repo *git.Repository) *tagsCommand {
	t := &tagsCommand{
		repo: repo,
	}
	cmd := &cobra.Command{
		Use:  "tags",
		RunE: t.runE,
	}
	t.cmd = cmd
	return t
}

func (t *tagsCommand) runE(cmd *cobra.Command, args []string) error {
	// get the tag iterator
	titr, err := t.repo.Tags()
	if err != nil {
		return err
	}
	versions := []*semver.Version{}
	err = titr.ForEach(func(ref *plumbing.Reference) error {
		ver, err := semver.NewVersion(ref.Name().Short())
		if err == nil && ver.Prerelease() == "" {
			versions = append(versions, ver)
		}
		return nil
	})
	sort.Slice(versions, func(i, j int) bool {
		v1 := versions[i]
		v2 := versions[j]
		return v2.LessThan(v1)
	})
	for _, ver := range versions {
		cmd.Println(ver.String())
	}
	return err
}

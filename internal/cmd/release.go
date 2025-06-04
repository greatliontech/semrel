package cmd

import (
	"fmt"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/greatliontech/semrel/internal/repository"
	"github.com/greatliontech/semrel/pkg/semrel"
	"github.com/spf13/cobra"
)

type releaseCommand struct {
	cmd               *cobra.Command
	repo              *repository.Repo
	cfg               *semrel.Config
	prerelease        string
	build             string
	currentBranchOnly bool
}

func newReleaseCommand(repo *repository.Repo, cfg *semrel.Config) *releaseCommand {
	c := &releaseCommand{
		repo: repo,
		cfg:  cfg,
	}
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Release a new version",
		RunE:  c.runE,
	}
	cmd.Flags().StringVarP(&c.prerelease, "prerelease", "p", "", "prerelease version")
	cmd.Flags().StringVarP(&c.build, "build", "b", "", "build version")
	cmd.Flags().BoolVarP(&c.currentBranchOnly, "current-branch-only", "", false, "only tags from the current branch")
	c.cmd = cmd
	return c
}

func (r *releaseCommand) runE(cmd *cobra.Command, args []string) error {
	var next semver.Version
	// check for initial version
	if r.cfg.InitialVersion() != nil {
		next = *r.cfg.InitialVersion()
	}

	// get latest tag version
	current, ref, err := r.repo.CurrentVersion(r.currentBranchOnly)
	if err != nil {
		return err
	}

	if !current.Equal(emptyVersion) {
		commits := []*semrel.Commit{}
		if ref != nil {
			commits, err = r.repo.Commits(plumbing.ZeroHash, ref.Hash())
			if err != nil {
				return err
			}
		}
		next = semrel.NextVersion(current, commits, r.cfg)
	}

	if next.Equal(current) {
		currentTag := fmt.Sprintf("%s%s", r.cfg.Prefix(), current.String())
		fmt.Println(currentTag)
		return nil
	}

	if r.prerelease != "" {
		next, err = next.SetPrerelease(r.prerelease)
		if err != nil {
			return err
		}
	}

	if r.build != "" {
		next, err = next.SetMetadata(r.build)
		if err != nil {
			return err
		}
	}

	nextTag := fmt.Sprintf("%s%s", r.cfg.Prefix(), next.String())
}

// matchRule holds one regex + replacement template.
type matchRule struct {
	Match   *regexp.Regexp
	Replace string
}

// Apply runs the regex on s, substituting every match
// according to the Replace template (using $1, $2, … for capture groups).
func (r *matchRule) Apply(s string) string {
	return r.Match.ReplaceAllString(s, r.Replace)
}

func generateReleaseNotes(commits []*semrel.Commit, cfg *semrel.Config) string {
	// This function should generate release notes based on the commits and configuration.
	// For simplicity, we will return a placeholder string.
	return "Release notes generated based on commits."
}

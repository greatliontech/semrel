package cmd

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/greatliontech/semrel/internal/release"
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

	commits := []*semrel.Commit{}
	if !current.Equal(emptyVersion) {
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

	var filters *release.Filters
	if r.cfg.Filters() != nil {
		filters = &release.Filters{}
		filters.Types = r.cfg.Filters().Types
		filters.Scopes = r.cfg.Filters().Scopes
	}

	rules := []*release.MatchRule{}
	if len(r.cfg.MatchRules()) > 0 {
		for _, rule := range r.cfg.MatchRules() {
			r, err := release.NewMatchRule(rule.Match, rule.Replace)
			if err != nil {
				return fmt.Errorf("invalid match rule: %w", err)
			}
			rules = append(rules, r)
		}
	}

	notes := release.GenerateReleaseNotes(commits, filters, rules)

	platform := r.cfg.Platform()
	if platform == "" {
		platform = os.Getenv("SEMREL_PLATFORM")
	}
	if platform == "" {
		platform, err = release.DetectPlatform()
		if err != nil {
			return err
		}
	}

	tok := os.Getenv("SEMREL_TOKEN")
	proj := os.Getenv("SEMREL_PROJECT")
	branch := os.Getenv("SEMREL_BRANCH")

	releaser, err := release.Platform(platform, tok, proj, branch)
	if err != nil {
		return err
	}

	if err := releaser.Release(nextTag, notes); err != nil {
		return fmt.Errorf("could not create release for next %q (current %q): %w", nextTag, current.String(), err)
	}

	fmt.Println(nextTag)
	return nil
}

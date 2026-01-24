package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
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
	root              *rootCommand
	prerelease        string
	build             string
	currentBranchOnly bool
	jsonOutput        bool
	notesOnly         bool
}

type ReleaseOutput struct {
	Current string `json:"current"`
	Next    string `json:"next"`
	Created bool   `json:"created"`
	DryRun  bool   `json:"dry_run,omitempty"`
}

func newReleaseCommand(repo *repository.Repo, cfg *semrel.Config, root *rootCommand) *releaseCommand {
	c := &releaseCommand{
		repo: repo,
		cfg:  cfg,
		root: root,
	}
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Release a new version",
		RunE:  c.runE,
	}
	cmd.Flags().StringVarP(&c.prerelease, "prerelease", "p", "", "prerelease version")
	cmd.Flags().StringVarP(&c.build, "build", "b", "", "build version")
	cmd.Flags().BoolVarP(&c.currentBranchOnly, "current-branch-only", "", false, "only tags from the current branch")
	cmd.Flags().BoolVar(&c.jsonOutput, "json", false, "output in JSON format")
	cmd.Flags().BoolVar(&c.notesOnly, "notes-only", false, "output release notes only, do not create release")
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

	currentTag := fmt.Sprintf("%s%s", r.cfg.Prefix(), current.String())

	if next.Equal(current) {
		if r.jsonOutput {
			return r.outputJSON(ReleaseOutput{
				Current: currentTag,
				Next:    currentTag,
				Created: false,
				DryRun:  r.root.dryRun,
			})
		}
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

	// Handle dry-run mode
	if r.root.dryRun {
		slog.Debug("dry-run: would create release", "tag", nextTag, "current", currentTag)
		if r.jsonOutput {
			return r.outputJSON(ReleaseOutput{
				Current: currentTag,
				Next:    nextTag,
				Created: false,
				DryRun:  true,
			})
		}
		fmt.Printf("[dry-run] would release: %s\n", nextTag)
		return nil
	}

	var filters *release.Filters
	if r.cfg.Filters() != nil {
		filters = &release.Filters{}
		filters.Types = r.cfg.Filters().Types
		filters.Scopes = r.cfg.Filters().Scopes
	}

	rules := []*release.MatchRule{}
	if len(r.cfg.MatchRules()) > 0 {
		for _, rule := range r.cfg.MatchRules() {
			mr, err := release.NewMatchRule(rule.Match, rule.Replace)
			if err != nil {
				return fmt.Errorf("invalid match rule: %w", err)
			}
			rules = append(rules, mr)
		}
	}

	notes := release.GenerateReleaseNotes(commits, filters, rules)

	// Handle notes-only mode
	if r.notesOnly {
		fmt.Print(notes)
		return nil
	}

	platform, tok, proj, err := release.DetectPlatform()
	// only one type of error here, release.ErrPlatformDetectionFailed
	if err != nil {
		platform = r.cfg.Platform()
		// env has precedence over config
		if p := os.Getenv("SEMREL_PLATFORM"); p != "" {
			platform = p
		}
		if platform == "" {
			return fmt.Errorf("platform not specified and could not be detected, please set SEMREL_PLATFORM environment variable or configure platform in semrel config: %w", err)
		}
	}

	// check for overrides from env
	if t := os.Getenv("SEMREL_TOKEN"); t != "" {
		tok = t
	}
	if p := os.Getenv("SEMREL_PROJECT"); p != "" {
		proj = p
	}
	// branch is explicitly set or empty
	branch := os.Getenv("SEMREL_BRANCH")

	releaser, err := release.Platform(platform, tok, proj, branch)
	if err != nil {
		return err
	}

	if err := releaser.Release(nextTag, notes); err != nil {
		return fmt.Errorf("could not create release for next %q (current %q): %w", nextTag, current.String(), err)
	}

	if r.jsonOutput {
		return r.outputJSON(ReleaseOutput{
			Current: currentTag,
			Next:    nextTag,
			Created: true,
		})
	}

	fmt.Println(nextTag)
	return nil
}

func (r *releaseCommand) outputJSON(output ReleaseOutput) error {
	b, err := json.Marshal(output)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

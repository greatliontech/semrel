package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/greatliontech/semrel/internal/errors"
	"github.com/greatliontech/semrel/internal/repository"
	"github.com/greatliontech/semrel/pkg/semrel"
	"github.com/spf13/cobra"
)

type rootCommand struct {
	cmd               *cobra.Command
	repo              *repository.Repo
	cfg               *semrel.Config
	currentBranchOnly bool
	createTag         bool
	pushTag           bool
	authUsername      string
	authPassword      string
	authToken         string
	prerelease        string
	build             string
	verbose           bool
	dryRun            bool
}

func New(rp *repository.Repo, cfg *semrel.Config, ver string) (*rootCommand, error) {
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
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			c.setupLogging()
		},
	}
	cmd.PersistentFlags().BoolVarP(&c.verbose, "verbose", "v", false, "enable verbose output")
	cmd.PersistentFlags().BoolVar(&c.dryRun, "dry-run", false, "simulate actions without making changes")
	cmd.Flags().StringVarP(&c.prerelease, "prerelease", "p", "", "prerelease version")
	cmd.Flags().StringVarP(&c.build, "build", "b", "", "build version")
	cmd.Flags().BoolVarP(&c.currentBranchOnly, "current-branch-only", "", false, "only tags from the current branch")
	cmd.Flags().BoolVarP(&c.createTag, "create-tag", "", false, "create the tag")
	cmd.Flags().BoolVarP(&c.pushTag, "push-tag", "", false, "push the tag")
	cmd.Flags().StringVarP(&c.authUsername, "auth-username", "", "", "username for basic auth")
	cmd.Flags().StringVarP(&c.authPassword, "auth-password", "", "", "password for basic auth")
	cmd.MarkFlagsRequiredTogether("auth-username", "auth-password")
	cmd.Flags().StringVarP(&c.authToken, "auth-token", "", "", "token for auth")
	cmd.MarkFlagsMutuallyExclusive("auth-username", "auth-token")
	cmd.MarkFlagsMutuallyExclusive("auth-password", "auth-token")
	cmd.AddCommand(
		newCurrentCommand(rp, cfg).cmd,
		newCompareCommand(rp, cfg).cmd,
		newValidateCommand().cmd,
		newReleaseCommand(rp, cfg, c).cmd,
		newMCPCommand(rp, cfg, ver).cmd,
		newInitCommand(rp).cmd,
	)
	c.cmd = cmd
	return c, nil
}

func (r *rootCommand) setupLogging() {
	level := slog.LevelWarn
	if r.verbose {
		level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))
}

func (r *rootCommand) Execute() int {
	err := r.cmd.Execute()
	if err != nil {
		if err != errCompareFailed {
			slog.Error("command failed", "error", err)
		}
		return exitCodeFromError(err)
	}
	return ExitSuccess
}

func exitCodeFromError(err error) int {
	switch {
	case err == errCompareFailed:
		return ExitGeneralError
	case errors.IsValidationError(err):
		return ExitValidation
	case errors.IsConfigError(err):
		return ExitConfig
	case errors.IsRepositoryError(err):
		return ExitRepository
	case errors.IsOperationError(err):
		return ExitRelease
	default:
		return ExitGeneralError
	}
}

var emptyVersion = semver.New(0, 0, 0, "", "")

func (r *rootCommand) runE(cmd *cobra.Command, args []string) error {
	var next semver.Version
	// check for initial version
	if r.cfg.InitialVersion() != nil {
		next = *r.cfg.InitialVersion()
	}

	// get latest tag version
	current, ref, err := r.repo.CurrentVersion(r.currentBranchOnly, r.cfg.StableOnly())
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

	if r.createTag || r.cfg.CreateTag() {
		head, err := r.repo.Head()
		if err != nil {
			return err
		}

		var auth transport.AuthMethod
		if un := os.Getenv("SEMREL_AUTH_USERNAME"); un != "" {
			r.authUsername = un
		}
		if pw := os.Getenv("SEMREL_AUTH_PASSWORD"); pw != "" {
			r.authPassword = pw
		}
		if tok := os.Getenv("SEMREL_AUTH_TOKEN"); tok != "" {
			r.authToken = tok
		}

		if r.authUsername != "" && r.authPassword != "" {
			auth = &http.BasicAuth{
				Username: r.authUsername,
				Password: r.authPassword,
			}
		}

		if r.authToken != "" {
			auth = &http.BasicAuth{
				Username: "git",
				Password: r.authToken,
			}
		}

		pushTag := r.pushTag || r.cfg.PushTag()
		err = r.repo.CreateTag(nextTag, head, pushTag, auth)
		if err != nil {
			return err
		}
	}

	fmt.Println(nextTag)
	return nil
}

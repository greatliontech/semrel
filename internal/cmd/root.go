package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
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
}

func New(r *git.Repository, cfg *semrel.Config, ver string) (*rootCommand, error) {
	rp := repository.New(r)
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
	}
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
		newCurrentCommand(rp).cmd,
		newCompareCommand(rp).cmd,
		newValidateCommand().cmd,
	)
	c.cmd = cmd
	return c, nil
}

func (r *rootCommand) Execute() {
	err := r.cmd.Execute()
	if err != nil {
		if err != errCompareFailed {
			slog.Error("command failed", "error", err)
		}
		os.Exit(1)
	}
	os.Exit(0)
}

func (r *rootCommand) runE(cmd *cobra.Command, args []string) error {
	// get latest tag version
	ver, ref, err := r.repo.CurrentVersion(r.currentBranchOnly)
	if err != nil {
		if err != repository.ErrNoTags {
			return err
		}
		ver = r.cfg.InitialVersion()
	}

	commits := []*semrel.Commit{}
	if ref != nil {
		commits, err = r.repo.Commits(plumbing.ZeroHash, ref.Hash())
		if err != nil {
			return err
		}
	}

	next := semrel.NextVersion(ver, commits, r.cfg)

	if next.Equal(ver) {
		fmt.Println(ver.String())
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

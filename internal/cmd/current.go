package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/greatliontech/semrel/internal/repository"
	"github.com/greatliontech/semrel/pkg/semrel"
	"github.com/spf13/cobra"
)

type currentCommand struct {
	cmd               *cobra.Command
	repo              *repository.Repo
	cfg               *semrel.Config
	currentBranchOnly bool
	jsonOutput        bool
}

type CurrentOutput struct {
	Version string `json:"version"`
	Tag     string `json:"tag"`
}

func newCurrentCommand(repo *repository.Repo, cfg *semrel.Config) *currentCommand {
	c := &currentCommand{
		repo: repo,
		cfg:  cfg,
	}
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Print the current release version",
		RunE:  c.runE,
	}
	cmd.Flags().BoolVarP(&c.currentBranchOnly, "current-branch-only", "", false, "only tags from the current branch")
	cmd.Flags().BoolVar(&c.jsonOutput, "json", false, "output in JSON format")
	c.cmd = cmd
	return c
}

func (c *currentCommand) runE(cmd *cobra.Command, args []string) error {
	cv, _, err := c.repo.CurrentVersion(c.currentBranchOnly)
	if err != nil {
		return err
	}
	currentTag := fmt.Sprintf("%s%s", c.cfg.Prefix(), cv.String())

	if c.jsonOutput {
		output := CurrentOutput{
			Version: cv.String(),
			Tag:     currentTag,
		}
		b, err := json.Marshal(output)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	}

	fmt.Println(currentTag)
	return nil
}

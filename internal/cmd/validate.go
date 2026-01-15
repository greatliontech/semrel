package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
)

type validateCommand struct {
	cmd          *cobra.Command
	strict       bool
	noPreRelease bool
	noBuild      bool
	jsonOutput   bool
}

type ValidateOutput struct {
	Input   string `json:"input"`
	Valid   bool   `json:"valid"`
	Version string `json:"version,omitempty"`
	Error   string `json:"error,omitempty"`
}

func newValidateCommand() *validateCommand {
	c := &validateCommand{}
	c.cmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate a semver version string",
		RunE:  c.runE,
		Args:  cobra.ExactArgs(1),
	}
	c.cmd.Flags().BoolVar(&c.strict, "strict", false, "strict semver validation")
	c.cmd.Flags().BoolVar(&c.noPreRelease, "noPrerelease", false, "do not allow pre-release versions")
	c.cmd.Flags().BoolVar(&c.noBuild, "noBuild", false, "do not allow build metadata")
	c.cmd.Flags().BoolVar(&c.jsonOutput, "json", false, "output in JSON format")
	return c
}

var errInvalidVersion = errors.New("invalid version")

func (c *validateCommand) runE(cmd *cobra.Command, args []string) error {
	vs := args[0]
	var sv *semver.Version
	var err error
	if c.strict {
		sv, err = semver.StrictNewVersion(vs)
	} else {
		sv, err = semver.NewVersion(vs)
	}
	if err != nil {
		if c.jsonOutput {
			return c.outputJSON(ValidateOutput{
				Input: vs,
				Valid: false,
				Error: err.Error(),
			})
		}
		return err
	}
	if c.noPreRelease && len(sv.Prerelease()) > 0 {
		if c.jsonOutput {
			return c.outputJSON(ValidateOutput{
				Input: vs,
				Valid: false,
				Error: "prerelease not allowed",
			})
		}
		return errInvalidVersion
	}
	if c.noBuild && len(sv.Metadata()) > 0 {
		if c.jsonOutput {
			return c.outputJSON(ValidateOutput{
				Input: vs,
				Valid: false,
				Error: "build metadata not allowed",
			})
		}
		return errInvalidVersion
	}

	if c.jsonOutput {
		return c.outputJSON(ValidateOutput{
			Input:   vs,
			Valid:   true,
			Version: sv.String(),
		})
	}

	return nil
}

func (c *validateCommand) outputJSON(output ValidateOutput) error {
	b, err := json.Marshal(output)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	if !output.Valid {
		return errInvalidVersion
	}
	return nil
}

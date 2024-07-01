package cmd

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
)

type validateCommand struct {
	cmd          *cobra.Command
	noPreRelease bool
	noBuild      bool
}

func newValidateCommand() *validateCommand {
	c := &validateCommand{}
	c.cmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate a semver version string",
		RunE:  c.runE,
		Args:  cobra.ExactArgs(1),
	}
	c.cmd.Flags().BoolVar(&c.noPreRelease, "noPrerelease", false, "do not allow pre-release versions")
	c.cmd.Flags().BoolVar(&c.noBuild, "noBuild", false, "do not allow build metadata")
	return c
}

var errInvalidVersion = errors.New("invalid version")

func (c *validateCommand) runE(cmd *cobra.Command, args []string) error {
	vs := args[0]
	sv, err := semver.NewVersion(vs)
	if err != nil {
		return err
	}
	if c.noPreRelease && len(sv.Prerelease()) > 0 {
		return errInvalidVersion
	}
	if c.noBuild && len(sv.Metadata()) > 0 {
		return errInvalidVersion
	}
	return nil
}

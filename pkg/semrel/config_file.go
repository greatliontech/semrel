package semrel

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Filters struct {
	Types  []string `yaml:"types"`
	Scopes []string `yaml:"scopes"`
}

type MatchRule struct {
	Match   string `yaml:"match"`
	Replace string `yaml:"replace"`
}

// ConfigFile is the configuration file for the semantic release tool in YAML format
type ConfigFile struct {
	// The default bump type if no commit types match. Default is "none"
	DefaultBump string `yaml:"defaultBump" json:"defaultBump" enum:"none,patch,minor,major" default:"none"`

	// The commit types that are considered as patch
	PatchTypes []string `yaml:"patchTypes" json:"patchTypes" default:"[\"fix\"]"`

	// The commit types that are considered as minor
	MinorTypes []string `yaml:"minorTypes" json:"minorTypes" default:"[\"feat\"]"`

	// The commit types that are considered as major
	MajorTypes []string `yaml:"majorTypes" json:"majorTypes"`

	// Initial version overrides the default initial version. If not set, it defaults to 1.0.0
	InitialVersion string `yaml:"initialVersion" json:"initialVersion" default:"1.0.0"`

	// Development if true, sets the initial version to 0.1.0 and treats breaking changes as minor bumps
	Development bool `yaml:"development" json:"development"`

	// Prefix is the prefix for the versions
	Prefix string `yaml:"prefix" json:"prefix"`

	// CreateTag if true, creates the next version tag
	CreateTag bool `yaml:"createTag"`

	// PushTag if true, pushes the next version tag. Requires CreateTag to be true
	PushTag bool `yaml:"pushTag"`

	// Platform that the tool is running on, e.g., "github", "gitlab", etc.
	Platform string `yaml:"platform"`

	// MatchRules are regex rules for matching commit messages and replacing them
	MatchRules []MatchRule `yaml:"matchRules"`

	// Filters are used to exclude certain commit types and scopes from release notes
	Filters *Filters `yaml:"filters"`
}

func ConfigFileFromPath(path string) (*ConfigFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ConfigFileFromBytes(b)
}

func ConfigFileFromBytes(b []byte) (*ConfigFile, error) {
	c := &ConfigFile{}
	err := yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

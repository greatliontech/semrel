package semrel

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ConfigFile is the configuration file for the semantic release tool in YAML format
type ConfigFile struct {
	// The commit types that are considered as no change
	DefaultBump string `yaml:"defaultBump"`

	// The commit types that are considered as patch
	PatchTypes []string `yaml:"patchTypes"`

	// The commit types that are considered as minor
	MinorTypes []string `yaml:"minorTypes"`

	// The commit types that are considered as major
	MajorTypes []string `yaml:"majorTypes"`

	// Initial version overrides the default initial version. If not set, it defaults to 1.0.0
	InitialVersion string `yaml:"initialVersion"`

	// Development if true, sets the initial version to 0.1.0 and treats breaking changes as minor bumps
	Development bool `yaml:"development"`

	// Prefix is the prefix for the versions
	Prefix string `yaml:"prefix"`

	// CreateTag if true, creates the next version tag
	CreateTag bool `yaml:"createTag"`

	// PushTag if true, pushes the next version tag. Requires CreateTag to be true
	PushTag bool `yaml:"pushTag"`
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

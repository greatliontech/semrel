package semrel

import (
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var ErrInvalidDefaultBump = errors.New("invalid default bump")

const (
	BumpPatch = "patch"
	BumpMinor = "minor"
	BumpMajor = "major"
	BumpNone  = "none"
)

type Config struct {
	// The commit types that are considered as patch
	PatchTypes []string `yaml:"patchTypes"`
	// The commit types that are considered as minor
	MinorTypes []string `yaml:"minorTypes"`
	// The commit types that are considered as major
	MajorTypes []string `yaml:"majorTypes"`
	// The commit types that are considered as no change
	NoChangeTypes []string `yaml:"noChangeTypes"`
	// The default bump if no changes are found
	DefaultBump string `yaml:"defaultBump"`
	// Whether the commit type is case sensitive
	CaseSensitive bool `yaml:"caseSensitive"`
}

// ParseConfig parses a YAML config file into a Config struct
func ParseConfig(b []byte) (*Config, error) {
	c := &Config{}
	err := yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	if c.DefaultBump == "" {
		c.DefaultBump = "none"
	}
	c.DefaultBump = strings.ToLower(c.DefaultBump)
	validDefaultBumps := []string{"patch", "minor", "major", "none"}
	isValidDefaultBump := false
	for _, bump := range validDefaultBumps {
		if c.DefaultBump == bump {
			isValidDefaultBump = true
			break
		}
	}
	if !isValidDefaultBump {
		return nil, ErrInvalidDefaultBump
	}
	return c, nil
}

// ParseConfigFile parses a YAML config file into a Config struct
func ParseConfigFile(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseConfig(b)
}

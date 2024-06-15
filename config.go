package semrel

import (
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var ErrInvalidDefaultBump = errors.New("invalid default bump")

var DefaultConfig = &Config{
	DefaultBump: BumpNone,
	PatchTypes:  []string{"fix"},
	MinorTypes:  []string{"feat"},
}

type BumpKind string

const (
	BumpNone  BumpKind = "none"
	BumpPatch BumpKind = "patch"
	BumpMinor BumpKind = "minor"
	BumpMajor BumpKind = "major"
)

func (b BumpKind) String() string {
	return string(b)
}

func (b BumpKind) IsValid() bool {
	if b == BumpNone || b == BumpPatch || b == BumpMinor || b == BumpMajor {
		return true
	}
	return false
}

func (b BumpKind) IsNone() bool {
	return b == BumpNone
}

func (b BumpKind) IsPatch() bool {
	return b == BumpPatch
}

func (b BumpKind) IsMinor() bool {
	return b == BumpMinor
}

func (b BumpKind) IsMajor() bool {
	return b == BumpMajor
}

func (b BumpKind) Index() int {
	switch b {
	case BumpNone:
		return 0
	case BumpPatch:
		return 1
	case BumpMinor:
		return 2
	case BumpMajor:
		return 3
	default:
		return -1
	}
}

type Config struct {
	// The commit types that are considered as no change
	DefaultBump BumpKind `yaml:"defaultBump"`
	// The commit types that are considered as patch
	PatchTypes []string `yaml:"patchTypes"`
	// The commit types that are considered as minor
	MinorTypes []string `yaml:"minorTypes"`
	// The commit types that are considered as major
	MajorTypes []string `yaml:"majorTypes"`
}

// ParseConfig parses a YAML config file into a Config struct
func ParseConfig(b []byte) (*Config, error) {
	c := &Config{}
	err := yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	if c.DefaultBump == "" {
		c.DefaultBump = BumpNone
	}
	c.DefaultBump = BumpKind(strings.ToLower(c.DefaultBump.String()))
	if !c.DefaultBump.IsValid() {
		return nil, ErrInvalidDefaultBump
	}
	if c.PatchTypes == nil {
		c.PatchTypes = []string{"fix"}
	}
	if c.MinorTypes == nil {
		c.MinorTypes = []string{"feat"}
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

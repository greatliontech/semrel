package semrel

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	mapset "github.com/deckarep/golang-set/v2"
)

type ConfigOption func(*Config)

var DefaultConfig = &Config{
	defaultBump:    BumpNone,
	patchTypes:     mapset.NewSet("fix"),
	minorTypes:     mapset.NewSet("feat"),
	majorTypes:     mapset.NewSet[string](),
	initialVersion: semver.New(1, 0, 0, "", ""),
}

func WithDefaultBump(b BumpKind) ConfigOption {
	return func(c *Config) {
		c.defaultBump = b
	}
}

func WithPatchTypes(types ...string) ConfigOption {
	return func(c *Config) {
		c.patchTypes = mapset.NewSet(types...)
	}
}

func WithMinorTypes(types ...string) ConfigOption {
	return func(c *Config) {
		c.minorTypes = mapset.NewSet(types...)
	}
}

func WithMajorTypes(types ...string) ConfigOption {
	return func(c *Config) {
		c.majorTypes = mapset.NewSet(types...)
	}
}

func WithInitialVersion(v *semver.Version) ConfigOption {
	return func(c *Config) {
		c.initialVersion = v
	}
}

func WithDevelopment() ConfigOption {
	return func(c *Config) {
		c.development = true
	}
}

type Config struct {
	defaultBump    BumpKind
	patchTypes     mapset.Set[string]
	minorTypes     mapset.Set[string]
	majorTypes     mapset.Set[string]
	initialVersion *semver.Version
	development    bool
}

func (c *Config) DefaultBump() BumpKind {
	return c.defaultBump
}

func (c *Config) BumpKind(str string) BumpKind {
	if c.patchTypes.Contains(str) {
		return BumpPatch
	}
	if c.minorTypes.Contains(str) {
		return BumpMinor
	}
	if c.majorTypes.Contains(str) {
		return BumpMajor
	}
	return BumpNone
}

func (c *Config) InitialVersion() *semver.Version {
	return c.initialVersion
}

func NewConfig(opts ...ConfigOption) (*Config, error) {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	if c.patchTypes.ContainsAny(c.minorTypes.ToSlice()...) ||
		c.patchTypes.ContainsAny(c.majorTypes.ToSlice()...) ||
		c.minorTypes.ContainsAny(c.majorTypes.ToSlice()...) {
		return nil, errors.New("commit types overlap")
	}
	if c.initialVersion == nil {
		if c.development {
			c.initialVersion = semver.New(0, 1, 0, "", "")
		} else {
			c.initialVersion = semver.New(1, 0, 0, "", "")
		}
	}
	return c, nil
}

func NewConfigFromConfigFile(cf *ConfigFile) (*Config, error) {
	opts := []ConfigOption{}

	if defBump := cf.DefaultBump; defBump != "" {
		bump, err := NewBump(defBump)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithDefaultBump(bump))
	}

	if len(cf.PatchTypes) > 0 {
		opts = append(opts, WithPatchTypes(cf.PatchTypes...))
	}

	if len(cf.MinorTypes) > 0 {
		opts = append(opts, WithMinorTypes(cf.MinorTypes...))
	}

	if len(cf.MajorTypes) > 0 {
		opts = append(opts, WithMajorTypes(cf.MajorTypes...))
	}

	if cf.InitialVersion != "" {
		v, err := semver.NewVersion(cf.InitialVersion)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithInitialVersion(v))
	}

	if cf.Development {
		opts = append(opts, WithDevelopment())
	}

	return NewConfig(opts...)
}

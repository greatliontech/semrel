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
	devMajorBump:   BumpPatch,
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

func WithDevelopmentMajorBump(b BumpKind) ConfigOption {
	return func(c *Config) {
		c.devMajorBump = b
	}
}

func WithPrefix(prefix string) ConfigOption {
	return func(c *Config) {
		c.prefix = prefix
	}
}

func WithCreateTag() ConfigOption {
	return func(c *Config) {
		c.createTag = true
	}
}

func WithPushTag() ConfigOption {
	return func(c *Config) {
		c.pushTag = true
	}
}

func WithPlatform(platform string) ConfigOption {
	return func(c *Config) {
		c.platform = platform
	}
}

func WithMaTchRules(rules ...MatchRule) ConfigOption {
	return func(c *Config) {
		c.matchRules = rules
	}
}

func WithFilters(filters *Filters) ConfigOption {
	return func(c *Config) {
		c.filters = filters
	}
}

type Config struct {
	patchTypes     mapset.Set[string]
	minorTypes     mapset.Set[string]
	majorTypes     mapset.Set[string]
	initialVersion *semver.Version
	prefix         string
	defaultBump    BumpKind
	devMajorBump   BumpKind
	development    bool
	createTag      bool
	pushTag        bool
	platform       string
	matchRules     []MatchRule
	filters        *Filters
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

func (c *Config) IsDevelopment() bool {
	return c.development
}

func (c *Config) Prefix() string {
	return c.prefix
}

func (c *Config) CreateTag() bool {
	return c.createTag
}

func (c *Config) PushTag() bool {
	return c.pushTag
}

func (c *Config) Platform() string {
	return c.platform
}

func (c *Config) MatchRules() []MatchRule {
	return c.matchRules
}

func (c *Config) Filters() *Filters {
	return c.filters
}

func NewConfig(opts ...ConfigOption) (*Config, error) {
	c := &Config{
		patchTypes: mapset.NewSet[string](),
		minorTypes: mapset.NewSet[string](),
		majorTypes: mapset.NewSet[string](),
	}
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

	if cf.PatchTypes == nil {
		opts = append(opts, WithPatchTypes(DefaultConfig.patchTypes.ToSlice()...))
	} else if len(cf.PatchTypes) > 0 {
		opts = append(opts, WithPatchTypes(cf.PatchTypes...))
	}

	if cf.MinorTypes == nil {
		opts = append(opts, WithMinorTypes(DefaultConfig.minorTypes.ToSlice()...))
	} else if len(cf.MinorTypes) > 0 {
		opts = append(opts, WithMinorTypes(cf.MinorTypes...))
	}

	if cf.MajorTypes == nil {
		opts = append(opts, WithMajorTypes(DefaultConfig.majorTypes.ToSlice()...))
	} else if len(cf.MajorTypes) > 0 {
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

	if cf.Prefix != "" {
		opts = append(opts, WithPrefix(cf.Prefix))
	}

	if cf.CreateTag {
		opts = append(opts, WithCreateTag())
	}

	if cf.PushTag {
		opts = append(opts, WithPushTag())
	}

	if cf.Platform != "" {
		opts = append(opts, WithPlatform(cf.Platform))
	}

	if len(cf.MatchRules) > 0 {
		opts = append(opts, WithMaTchRules(cf.MatchRules...))
	}

	if cf.Filters != nil {
		opts = append(opts, WithFilters(cf.Filters))
	}

	return NewConfig(opts...)
}

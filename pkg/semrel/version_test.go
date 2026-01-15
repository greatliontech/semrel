package semrel

import (
	"testing"

	"github.com/Masterminds/semver/v3"
)

func TestNextVersion(t *testing.T) {
	tests := []struct {
		name     string
		current  *semver.Version
		commits  []*Commit
		cfg      *Config
		expected string
	}{
		{
			name:    "patch bump from fix commit",
			current: semver.MustParse("1.0.0"),
			commits: []*Commit{
				{Type: "fix", Description: "fix something"},
			},
			cfg:      mustConfig(t),
			expected: "1.0.1",
		},
		{
			name:    "minor bump from feat commit",
			current: semver.MustParse("1.0.0"),
			commits: []*Commit{
				{Type: "feat", Description: "add feature"},
			},
			cfg:      mustConfig(t),
			expected: "1.1.0",
		},
		{
			name:    "major bump from breaking change attention marker",
			current: semver.MustParse("1.0.0"),
			commits: []*Commit{
				{Type: "feat", Attention: true, Description: "breaking feature"},
			},
			cfg:      mustConfig(t),
			expected: "2.0.0",
		},
		{
			name:    "major bump from BREAKING CHANGE in body",
			current: semver.MustParse("1.0.0"),
			commits: []*Commit{
				{Type: "feat", Description: "add feature", Body: "BREAKING CHANGE: this breaks things"},
			},
			cfg:      mustConfig(t),
			expected: "2.0.0",
		},
		{
			name:    "highest bump wins - major over minor",
			current: semver.MustParse("1.0.0"),
			commits: []*Commit{
				{Type: "fix", Description: "fix something"},
				{Type: "feat", Description: "add feature"},
				{Type: "feat", Attention: true, Description: "breaking feature"},
			},
			cfg:      mustConfig(t),
			expected: "2.0.0",
		},
		{
			name:    "highest bump wins - minor over patch",
			current: semver.MustParse("1.0.0"),
			commits: []*Commit{
				{Type: "fix", Description: "fix something"},
				{Type: "feat", Description: "add feature"},
			},
			cfg:      mustConfig(t),
			expected: "1.1.0",
		},
		{
			name:    "empty commits with default bump none",
			current: semver.MustParse("1.0.0"),
			commits: []*Commit{},
			cfg:      mustConfig(t),
			expected: "1.0.0",
		},
		{
			name:    "empty commits with default bump patch",
			current: semver.MustParse("1.0.0"),
			commits: []*Commit{},
			cfg:      mustConfig(t, WithDefaultBump(BumpPatch)),
			expected: "1.0.1",
		},
		{
			name:    "unknown commit type uses default bump",
			current: semver.MustParse("1.0.0"),
			commits: []*Commit{
				{Type: "chore", Description: "cleanup"},
			},
			cfg:      mustConfig(t, WithDefaultBump(BumpPatch)),
			expected: "1.0.1",
		},
		{
			name:    "development mode - major=0 with breaking change bumps patch",
			current: semver.MustParse("0.5.0"),
			commits: []*Commit{
				{Type: "feat", Attention: true, Description: "breaking feature"},
			},
			cfg:      mustConfig(t, WithDevelopment()),
			expected: "0.5.1",
		},
		{
			name:    "development mode - major>0 with breaking change bumps major",
			current: semver.MustParse("1.0.0"),
			commits: []*Commit{
				{Type: "feat", Attention: true, Description: "breaking feature"},
			},
			cfg:      mustConfig(t, WithDevelopment()),
			expected: "2.0.0",
		},
		{
			name:    "development mode - minor bump unchanged",
			current: semver.MustParse("0.5.0"),
			commits: []*Commit{
				{Type: "feat", Description: "add feature"},
			},
			cfg:      mustConfig(t, WithDevelopment()),
			expected: "0.6.0",
		},
		{
			name:    "custom major types",
			current: semver.MustParse("1.0.0"),
			commits: []*Commit{
				{Type: "breaking", Description: "breaking change"},
			},
			cfg:      mustConfig(t, WithMajorTypes("breaking")),
			expected: "2.0.0",
		},
		{
			name:    "prerelease version bumps to release",
			current: semver.MustParse("1.0.0-alpha.1"),
			commits: []*Commit{
				{Type: "fix", Description: "fix something"},
			},
			cfg:      mustConfig(t),
			expected: "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NextVersion(tt.current, tt.commits, tt.cfg)
			if result.String() != tt.expected {
				t.Errorf("NextVersion() = %s, want %s", result.String(), tt.expected)
			}
		})
	}
}

func mustConfig(t *testing.T, opts ...ConfigOption) *Config {
	t.Helper()
	baseOpts := []ConfigOption{
		WithPatchTypes("fix"),
		WithMinorTypes("feat"),
	}
	cfg, err := NewConfig(append(baseOpts, opts...)...)
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}
	return cfg
}

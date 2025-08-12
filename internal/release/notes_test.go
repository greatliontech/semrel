package release

import (
	"regexp"
	"testing"

	"github.com/greatliontech/semrel/pkg/semrel"
)

var testCommits = []*semrel.Commit{
	{
		Type:        "feat",
		Scope:       "core",
		Description: "Add new feature [TRACK-123]",
	},
	{
		Type:        "docs",
		Scope:       "readme",
		Description: "Update README [TRACK-456]",
	},
	{
		Type:        "fix",
		Scope:       "core",
		Description: "Fix a bug [TRACK-789]",
	},
}

func TestNotesSimple(t *testing.T) {
	notes := GenerateReleaseNotes(testCommits, nil, nil)
	exp := `- feat(core): Add new feature [TRACK-123]
- docs(readme): Update README [TRACK-456]
- fix(core): Fix a bug [TRACK-789]
`

	if notes != exp {
		t.Errorf("expected:\n%s\ngot:\n%s", exp, notes)
	}
}

func TestNotesFilters(t *testing.T) {
	filters := &Filters{
		Types: []string{"docs"},
	}
	notes := GenerateReleaseNotes(testCommits, filters, nil)
	exp := `- feat(core): Add new feature [TRACK-123]
- fix(core): Fix a bug [TRACK-789]
`

	if notes != exp {
		t.Errorf("expected:\n%s\ngot:\n%s", exp, notes)
	}
}

func TestNotesMatchRules(t *testing.T) {
	rules := []*MatchRule{
		{
			Match:   regexp.MustCompile(`\[(TRACK-\d+)\]`),
			Replace: `[#$1](https://example.com/issue/$1)`,
		},
	}
	notes := GenerateReleaseNotes(testCommits, nil, rules)
	exp := `- feat(core): Add new feature [#TRACK-123](https://example.com/issue/TRACK-123)
- docs(readme): Update README [#TRACK-456](https://example.com/issue/TRACK-456)
- fix(core): Fix a bug [#TRACK-789](https://example.com/issue/TRACK-789)
`

	if notes != exp {
		t.Errorf("expected:\n%s\ngot:\n%s", exp, notes)
	}
}

func TestNotesNoScopes(t *testing.T) {
	commits := []*semrel.Commit{
		{
			Type:        "feat",
			Description: "Add new feature [TRACK-123]",
		},
		{
			Type:        "fix",
			Description: "Fix a bug [TRACK-789]",
		},
	}
	notes := GenerateReleaseNotes(commits, nil, nil)
	exp := `- feat: Add new feature [TRACK-123]
- fix: Fix a bug [TRACK-789]
`
	if notes != exp {
		t.Errorf("expected:\n%s\ngot:\n%s", exp, notes)
	}
}

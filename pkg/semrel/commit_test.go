package semrel

import "testing"

var validCommitWithTrailers = `fix(problem): this is a valid commit message

This is a 2
line section

Signed-off-by: The Grumpy Lion
Some-other-footer: Some value
`

func TestValidCommitWithTrailers(t *testing.T) {
	c, err := ParseCommitMessage(validCommitWithTrailers)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if c.Type != "fix" {
		t.Errorf("expected type to be 'fix', got %s", c.Type)
	}

	if c.Scope != "problem" {
		t.Errorf("expected scope to be 'problem', got %s", c.Scope)
	}

	if c.Attention != false {
		t.Errorf("expected attention to be false, got %t", c.Attention)
	}

	if c.Description != "this is a valid commit message" {
		t.Errorf("expected description to be 'this is a valid commit message', got %s", c.Description)
	}

	wantBody := `This is a 2
line section
`

	if c.Body != wantBody {
		t.Errorf("expected body to be %q, got %q", wantBody, c.Body)
	}

	if len(c.Footers) != 2 {
		t.Errorf("expected 2 footers, got %d", len(c.Footers))
	}

	if c.Footers["Signed-off-by"] != "The Grumpy Lion" {
		t.Errorf("expected 'Signed-off-by' to be 'The Grumpy Lion', got %s", c.Footers["Signed-off-by"])
	}

	if c.Footers["Some-other-footer"] != "Some value" {
		t.Errorf("expected 'Some-other-footer' to be 'Some value', got %s", c.Footers["Some-other-footer"])
	}
}

func TestBreakingChangeAttentionMarker(t *testing.T) {
	msg := "feat!: this is a breaking change"
	c, err := ParseCommitMessage(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Type != "feat" {
		t.Errorf("expected type 'feat', got %s", c.Type)
	}
	if !c.Attention {
		t.Error("expected attention to be true for breaking change marker")
	}
}

func TestBreakingChangeInBody(t *testing.T) {
	msg := `feat: add new feature

This adds a new feature.

BREAKING CHANGE: this removes old behavior
`
	c, err := ParseCommitMessage(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Type != "feat" {
		t.Errorf("expected type 'feat', got %s", c.Type)
	}
	cfg, _ := NewConfig(WithPatchTypes("fix"), WithMinorTypes("feat"))
	bump := c.BumpKind(cfg)
	if bump != BumpMajor {
		t.Errorf("expected BumpMajor for BREAKING CHANGE in body, got %s", bump)
	}
}

func TestBreakingChangesPlural(t *testing.T) {
	msg := `fix: fix something

BREAKING CHANGES: multiple things changed
`
	c, err := ParseCommitMessage(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cfg, _ := NewConfig(WithPatchTypes("fix"), WithMinorTypes("feat"))
	bump := c.BumpKind(cfg)
	if bump != BumpMajor {
		t.Errorf("expected BumpMajor for BREAKING CHANGES in body, got %s", bump)
	}
}

func TestInvalidCommitFormat(t *testing.T) {
	invalidMessages := []string{
		"this is not a conventional commit",
		"missing colon type description",
		"",
		"no-colon",
		"type without description:",
	}
	for _, msg := range invalidMessages {
		_, err := ParseCommitMessage(msg)
		if err != ErrNotConventionalCommit {
			t.Errorf("expected ErrNotConventionalCommit for %q, got %v", msg, err)
		}
	}
}

func TestCommitWithoutScope(t *testing.T) {
	msg := "fix: simple fix without scope"
	c, err := ParseCommitMessage(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Type != "fix" {
		t.Errorf("expected type 'fix', got %s", c.Type)
	}
	if c.Scope != "" {
		t.Errorf("expected empty scope, got %s", c.Scope)
	}
	if c.Description != "simple fix without scope" {
		t.Errorf("expected description 'simple fix without scope', got %s", c.Description)
	}
}

func TestCommitTypeCaseInsensitive(t *testing.T) {
	msg := "FIX: uppercase type"
	c, err := ParseCommitMessage(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Type != "fix" {
		t.Errorf("expected type 'fix' (lowercase), got %s", c.Type)
	}
}

func TestCommitBumpKind(t *testing.T) {
	cfg, _ := NewConfig(
		WithPatchTypes("fix", "chore"),
		WithMinorTypes("feat"),
		WithMajorTypes("breaking"),
	)
	tests := []struct {
		name     string
		commit   *Commit
		expected BumpKind
	}{
		{"fix type", &Commit{Type: "fix"}, BumpPatch},
		{"chore type", &Commit{Type: "chore"}, BumpPatch},
		{"feat type", &Commit{Type: "feat"}, BumpMinor},
		{"breaking type", &Commit{Type: "breaking"}, BumpMajor},
		{"attention marker", &Commit{Type: "feat", Attention: true}, BumpMajor},
		{"unknown type", &Commit{Type: "docs"}, BumpNone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.commit.BumpKind(cfg)
			if result != tt.expected {
				t.Errorf("BumpKind() = %s, want %s", result, tt.expected)
			}
		})
	}
}

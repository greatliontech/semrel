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

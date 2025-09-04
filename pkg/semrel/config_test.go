package semrel

import "testing"

func TestConfigFromFileWithEmptyTypes(t *testing.T) {
	cnf := `
patchTypes: []
minorTypes:
  - feat
  - minor
`
	cf, err := ConfigFileFromBytes([]byte(cnf))
	if err != nil {
		t.Fatal(err)
	}

	c, err := NewConfigFromConfigFile(cf)
	if err != nil {
		t.Fatal(err)
	}

	if c.BumpKind("fix") != BumpNone {
		t.Errorf("expected no bump for 'fix', got %s", c.BumpKind("fix"))
	}

	if c.BumpKind("feat") != BumpMinor {
		t.Errorf("expected minor bump for 'feat', got %s", c.BumpKind("feat"))
	}

	if c.BumpKind("breaking") != BumpNone {
		t.Errorf("expected no bump for 'breaking', got %s", c.BumpKind("breaking"))
	}
}

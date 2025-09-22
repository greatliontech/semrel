package semrel

import (
	"os"
	"path/filepath"
	"testing"
)

var typesString = `
patchTypes:
  - fix
  - chore
minorTypes:
  - feat
  - minor
majorTypes:
  - breaking
platform: github
matchRules:
  - match: ((?:PENG|FE)-\d+)
    replace: myjira.net/$1
`

func TestConfigFileTypes(t *testing.T) {
	expectedPatchTypes := []string{"fix", "chore"}
	expectedMinorTypes := []string{"feat", "minor"}
	expectedMajorTypes := []string{"breaking"}

	c, err := ConfigFileFromBytes([]byte(typesString))
	if err != nil {
		t.Fatal(err)
	}

	if !expectedTypesMatch(expectedPatchTypes, c.PatchTypes) {
		t.Errorf("expected %v, got %v", expectedPatchTypes, c.PatchTypes)
	}
	if !expectedTypesMatch(expectedMinorTypes, c.MinorTypes) {
		t.Errorf("expected %v, got %v", expectedMinorTypes, c.MinorTypes)
	}
	if !expectedTypesMatch(expectedMajorTypes, c.MajorTypes) {
		t.Errorf("expected %v, got %v", expectedMajorTypes, c.MajorTypes)
	}
	if c.Platform != "github" {
		t.Errorf("expected platform 'github', got '%s'", c.Platform)
	}
	if len(c.MatchRules) != 1 {
		t.Errorf("expected 1 match rule, got %d", len(c.MatchRules))
	}
	rule := c.MatchRules[0]
	if rule.Match != `((?:PENG|FE)-\d+)` || rule.Replace != `myjira.net/$1` {
		t.Errorf("expected match rule to be %s -> %s, got %s -> %s", `((?:PENG|FE)-\d+)`, `myjira.net/$1`, rule.Match, rule.Replace)
	}
	if rule.Replace != "myjira.net/$1" {
		t.Errorf("expected replace to be 'myjira.net/$1', got '%s'", rule.Replace)
	}
}

func TestMinorDefaultBump(t *testing.T) {
	s := `
defaultBump: minor
`
	c, err := ConfigFileFromBytes([]byte(s))
	if err != nil {
		t.Fatal(err)
	}
	if c.DefaultBump != BumpMinor.String() {
		t.Errorf("expected %s, got %s", BumpMinor, c.DefaultBump)
	}
}

func TestConfigWithTabIndentation(t *testing.T) {
	s := `
patchTypes:
	- fix
	- chore
`

	_, err := ConfigFileFromBytes([]byte(s))
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestConfigFromFile(t *testing.T) {
	cf := filepath.Join(os.TempDir(), ".semrel")
	err := os.WriteFile(cf, []byte(typesString), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	c, err := ConfigFileFromPath(cf)
	if err != nil {
		t.Fatal(err)
	}
	expectedPatchTypes := []string{"fix", "chore"}
	expectedMinorTypes := []string{"feat", "minor"}
	expectedMajorTypes := []string{"breaking"}
	if !expectedTypesMatch(expectedPatchTypes, c.PatchTypes) {
		t.Errorf("expected %v, got %v", expectedPatchTypes, c.PatchTypes)
	}
	if !expectedTypesMatch(expectedMinorTypes, c.MinorTypes) {
		t.Errorf("expected %v, got %v", expectedMinorTypes, c.MinorTypes)
	}
	if !expectedTypesMatch(expectedMajorTypes, c.MajorTypes) {
		t.Errorf("expected %v, got %v", expectedMajorTypes, c.MajorTypes)
	}
}

func TestInvalidFilePath(t *testing.T) {
	_, err := ConfigFileFromPath("invalid/file/path")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestEmptyTypes(t *testing.T) {
	cnf := `
patchTypes: []
minorTypes:
  - feat
  - minor
`
	c, err := ConfigFileFromBytes([]byte(cnf))
	if err != nil {
		t.Fatal(err)
	}

	if c.PatchTypes == nil {
		t.Error("expected patchTypes to be an empty slice, got nil")
	}
	if c.MajorTypes != nil {
		t.Errorf("expected majorTypes to be nil, got %v", c.MajorTypes)
	}
	if len(c.MinorTypes) != 2 {
		t.Errorf("expected minorTypes to have 2 elements, got %d", len(c.MinorTypes))
	}
	for i, expected := range []string{"feat", "minor"} {
		if c.MinorTypes[i] != expected {
			t.Errorf("expected minorTypes[%d] to be %s, got %s", i, expected, c.MinorTypes[i])
		}
	}
}

func TestGithubIssue3(t *testing.T) {
	cnf := `
patchTypes:
  - fix
  - chore
minorTypes:
  - feat
`
	c, err := ConfigFileFromBytes([]byte(cnf))
	if err != nil {
		t.Fatal(err)
	}
	if len(c.MajorTypes) != 0 {
		t.Errorf("expected majorTypes to be empty, got %v", c.MajorTypes)
	}
	if len(c.PatchTypes) != 2 {
		t.Errorf("expected patchTypes to have 2 elements, got %d", len(c.PatchTypes))
	}
	if len(c.MinorTypes) != 1 {
		t.Errorf("expected minorTypes to have 1 element, got %d", len(c.MinorTypes))
	}
	if c.MinorTypes[0] != "feat" {
		t.Errorf("expected minorTypes[0] to be 'feat', got '%s'", c.MinorTypes[0])
	}
	if c.PatchTypes[0] != "fix" || c.PatchTypes[1] != "chore" {
		t.Errorf("expected patchTypes to be ['fix', 'chore'], got %v", c.PatchTypes)
	}
}

func expectedTypesMatch(expected, actual []string) bool {
	if len(expected) != len(actual) {
		return false
	}
	for i, v := range expected {
		if v != actual[i] {
			return false
		}
	}
	return true
}

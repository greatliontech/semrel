package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary once before running tests
	tmpDir, err := os.MkdirTemp("", "semrel-test-bin")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	binaryPath = filepath.Join(tmpDir, "semrel")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic("failed to build semrel: " + err.Error())
	}

	os.Exit(m.Run())
}

// testRepo creates an isolated git repository for testing
type testRepo struct {
	dir string
	t   *testing.T
}

func newTestRepo(t *testing.T) *testRepo {
	t.Helper()
	dir, err := os.MkdirTemp("", "semrel-test-repo")
	if err != nil {
		t.Fatal(err)
	}

	tr := &testRepo{dir: dir, t: t}
	tr.git("init")
	tr.git("config", "user.email", "test@test.com")
	tr.git("config", "user.name", "Test")

	return tr
}

func (tr *testRepo) cleanup() {
	os.RemoveAll(tr.dir)
}

func (tr *testRepo) git(args ...string) string {
	tr.t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = tr.dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		tr.t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func (tr *testRepo) commit(msg string) {
	tr.t.Helper()
	// Create a unique file for each commit
	f, err := os.CreateTemp(tr.dir, "file-*.txt")
	if err != nil {
		tr.t.Fatal(err)
	}
	f.WriteString(msg)
	f.Close()

	tr.git("add", ".")
	tr.git("commit", "-m", msg)
}

func (tr *testRepo) tag(tag string) {
	tr.t.Helper()
	tr.git("tag", tag)
}

func (tr *testRepo) semrel(args ...string) (string, error) {
	tr.t.Helper()
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = tr.dir
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func (tr *testRepo) mustSemrel(args ...string) string {
	tr.t.Helper()
	out, err := tr.semrel(args...)
	if err != nil {
		tr.t.Fatalf("semrel %v failed: %v\n%s", args, err, out)
	}
	return out
}

// CLI Integration Tests

func TestCLI_CurrentVersion(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial commit")
	repo.tag("v1.0.0")

	out := repo.mustSemrel("current")
	if out != "1.0.0" {
		t.Errorf("current = %q, want %q", out, "1.0.0")
	}
}

func TestCLI_NextVersion_PatchBump(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")
	repo.tag("v1.0.0")
	repo.commit("fix: bug fix")

	out := repo.mustSemrel()
	if out != "1.0.1" {
		t.Errorf("next version = %q, want %q", out, "1.0.1")
	}
}

func TestCLI_NextVersion_MinorBump(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")
	repo.tag("v1.0.0")
	repo.commit("feat: new feature")

	out := repo.mustSemrel()
	if out != "1.1.0" {
		t.Errorf("next version = %q, want %q", out, "1.1.0")
	}
}

func TestCLI_NextVersion_MajorBump_AttentionMarker(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")
	repo.tag("v1.0.0")
	repo.commit("feat!: breaking change")

	out := repo.mustSemrel()
	if out != "2.0.0" {
		t.Errorf("next version = %q, want %q", out, "2.0.0")
	}
}

func TestCLI_NextVersion_MajorBump_BreakingChangeFooter(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")
	repo.tag("v1.0.0")
	repo.commit("feat: add feature\n\nBREAKING CHANGE: this breaks things")

	out := repo.mustSemrel()
	if out != "2.0.0" {
		t.Errorf("next version = %q, want %q", out, "2.0.0")
	}
}

func TestCLI_NextVersion_NonConventionalCommit_NoBump(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")
	repo.tag("v1.0.0")
	repo.commit("random commit message without type")

	out := repo.mustSemrel()
	if out != "1.0.0" {
		t.Errorf("next version = %q, want %q (no bump)", out, "1.0.0")
	}
}

func TestCLI_NextVersion_HighestBumpWins(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")
	repo.tag("v1.0.0")
	repo.commit("fix: bug fix")
	repo.commit("feat: new feature")
	repo.commit("fix: another fix")

	out := repo.mustSemrel()
	if out != "1.1.0" {
		t.Errorf("next version = %q, want %q (minor wins over patch)", out, "1.1.0")
	}
}

func TestCLI_NextVersion_FreshRepo(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial commit")

	out := repo.mustSemrel()
	// Default initial version is 1.0.0
	if out != "1.0.0" {
		t.Errorf("fresh repo version = %q, want %q", out, "1.0.0")
	}
}

func TestCLI_Prerelease(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")
	repo.tag("v1.0.0")
	repo.commit("fix: bug fix")

	out := repo.mustSemrel("--prerelease", "alpha.1")
	if !strings.HasSuffix(out, "-alpha.1") {
		t.Errorf("prerelease version = %q, want suffix %q", out, "-alpha.1")
	}
}

func TestCLI_BuildMetadata(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")
	repo.tag("v1.0.0")
	repo.commit("fix: bug fix")

	out := repo.mustSemrel("--build", "sha.abc123")
	if !strings.HasSuffix(out, "+sha.abc123") {
		t.Errorf("build version = %q, want suffix %q", out, "+sha.abc123")
	}
}

func TestCLI_Current_JSONOutput(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")
	repo.tag("v1.2.3")

	out := repo.mustSemrel("current", "--json")

	var result struct {
		Version string `json:"version"`
		Tag     string `json:"tag"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, out)
	}
	if result.Version != "1.2.3" {
		t.Errorf("version = %q, want %q", result.Version, "1.2.3")
	}
}

func TestCLI_Validate_Valid(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")

	_, err := repo.semrel("validate", "1.2.3")
	if err != nil {
		t.Errorf("validate 1.2.3 should succeed, got error: %v", err)
	}
}

func TestCLI_Validate_Invalid(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")

	_, err := repo.semrel("validate", "not-a-version")
	if err == nil {
		t.Error("validate not-a-version should fail")
	}
}

func TestCLI_Validate_JSONOutput(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")

	out := repo.mustSemrel("validate", "1.2.3", "--json")

	var result struct {
		Input   string `json:"input"`
		Valid   bool   `json:"valid"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, out)
	}
	if !result.Valid {
		t.Error("valid should be true")
	}
	if result.Version != "1.2.3" {
		t.Errorf("version = %q, want %q", result.Version, "1.2.3")
	}
}

func TestCLI_Compare_Success(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")

	tests := []struct {
		name string
		args []string
	}{
		{"ge satisfied", []string{"compare", "1.5.0", "--ge", "1.0.0"}},
		{"le satisfied", []string{"compare", "1.0.0", "--le", "1.5.0"}},
		{"gt satisfied", []string{"compare", "2.0.0", "--gt", "1.0.0"}},
		{"lt satisfied", []string{"compare", "1.0.0", "--lt", "2.0.0"}},
		{"equal with ge", []string{"compare", "1.0.0", "--ge", "1.0.0"}},
		{"equal with le", []string{"compare", "1.0.0", "--le", "1.0.0"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.semrel(tt.args...)
			if err != nil {
				t.Errorf("compare %v should succeed, got error", tt.args)
			}
		})
	}
}

func TestCLI_Compare_Failure(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")

	tests := []struct {
		name string
		args []string
	}{
		{"ge not satisfied", []string{"compare", "1.0.0", "--ge", "2.0.0"}},
		{"le not satisfied", []string{"compare", "2.0.0", "--le", "1.0.0"}},
		{"gt not satisfied", []string{"compare", "1.0.0", "--gt", "1.0.0"}},
		{"lt not satisfied", []string{"compare", "1.0.0", "--lt", "1.0.0"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.semrel(tt.args...)
			if err == nil {
				t.Errorf("compare %v should fail", tt.args)
			}
		})
	}
}

func TestCLI_Compare_JSONOutput(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")

	out := repo.mustSemrel("compare", "1.5.0", "--ge", "1.0.0", "--json")

	var result struct {
		Version string `json:"version"`
		Result  bool   `json:"result"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, out)
	}
	if !result.Result {
		t.Error("result should be true")
	}
}

func TestCLI_DryRun(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")
	repo.tag("v1.0.0")
	repo.commit("fix: bug fix")

	out := repo.mustSemrel("release", "--dry-run")
	if !strings.Contains(out, "dry-run") {
		t.Errorf("dry-run output should contain 'dry-run', got: %s", out)
	}
}

func TestCLI_DryRun_JSONOutput(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")
	repo.tag("v1.0.0")
	repo.commit("fix: bug fix")

	out := repo.mustSemrel("release", "--dry-run", "--json")

	var result struct {
		Current string `json:"current"`
		Next    string `json:"next"`
		Created bool   `json:"created"`
		DryRun  bool   `json:"dry_run"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, out)
	}
	if !result.DryRun {
		t.Error("dry_run should be true")
	}
	if result.Created {
		t.Error("created should be false in dry-run mode")
	}
}

func TestCLI_Verbose(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	repo.commit("feat: initial")
	repo.tag("v1.0.0")
	repo.commit("fix: bug fix")

	// Just verify it doesn't crash
	_, err := repo.semrel("--verbose")
	if err != nil {
		t.Errorf("--verbose should not cause error: %v", err)
	}
}

func TestCLI_ConfigFile(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	// Create a config file with custom settings
	configContent := `
patchTypes:
  - fix
  - chore
minorTypes:
  - feat
prefix: "v"
`
	configPath := filepath.Join(repo.dir, ".semrel.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	repo.commit("feat: initial")
	repo.tag("v1.0.0")
	repo.commit("chore: cleanup")

	out := repo.mustSemrel()
	// chore should trigger patch bump
	if out != "v1.0.1" {
		t.Errorf("version = %q, want %q (chore should be patch)", out, "v1.0.1")
	}
}

func TestCLI_ConfigFile_Development(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	// Create a config with development mode
	configContent := `
development: true
`
	configPath := filepath.Join(repo.dir, ".semrel.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	repo.commit("feat: initial")

	out := repo.mustSemrel()
	// Development mode should start at 0.1.0
	if out != "0.1.0" {
		t.Errorf("development initial version = %q, want %q", out, "0.1.0")
	}
}

func TestCLI_ConfigFile_Development_BreakingChange(t *testing.T) {
	repo := newTestRepo(t)
	defer repo.cleanup()

	configContent := `
development: true
`
	configPath := filepath.Join(repo.dir, ".semrel.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	repo.commit("feat: initial")
	repo.tag("v0.1.0")
	repo.commit("feat!: breaking change")

	out := repo.mustSemrel()
	// In development mode with major=0, breaking changes bump patch
	if out != "0.1.1" {
		t.Errorf("dev mode breaking change = %q, want %q", out, "0.1.1")
	}
}

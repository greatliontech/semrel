package release

import (
	"os"
	"testing"
)

func TestPlatformGitHub(t *testing.T) {
	r, err := Platform("github", "token", "owner/repo", "main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil releaser")
	}
}

func TestPlatformGitLab(t *testing.T) {
	r, err := Platform("gitlab", "token", "12345", "main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil releaser")
	}
}

func TestPlatformCaseInsensitive(t *testing.T) {
	platforms := []string{"GITHUB", "GitHub", "GiThUb"}
	for _, p := range platforms {
		r, err := Platform(p, "token", "owner/repo", "main")
		if err != nil {
			t.Errorf("Platform(%q) unexpected error: %v", p, err)
		}
		if r == nil {
			t.Errorf("Platform(%q) returned nil releaser", p)
		}
	}
}

func TestPlatformUnsupported(t *testing.T) {
	_, err := Platform("bitbucket", "token", "project", "main")
	if err == nil {
		t.Fatal("expected error for unsupported platform")
	}
	unsupportedErr, ok := err.(*ErrUnsupportedPlatform)
	if !ok {
		t.Fatalf("expected ErrUnsupportedPlatform, got %T", err)
	}
	if unsupportedErr.platform != "bitbucket" {
		t.Errorf("expected platform 'bitbucket', got %s", unsupportedErr.platform)
	}
}

func TestDetectPlatformGitHub(t *testing.T) {
	os.Setenv("GITHUB_ACTIONS", "true")
	os.Setenv("GITHUB_TOKEN", "gh-token")
	os.Setenv("GITHUB_REPOSITORY", "owner/repo")
	defer func() {
		os.Unsetenv("GITHUB_ACTIONS")
		os.Unsetenv("GITHUB_TOKEN")
		os.Unsetenv("GITHUB_REPOSITORY")
	}()

	platform, token, project, err := DetectPlatform()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if platform != "github" {
		t.Errorf("expected platform 'github', got %s", platform)
	}
	if token != "gh-token" {
		t.Errorf("expected token 'gh-token', got %s", token)
	}
	if project != "owner/repo" {
		t.Errorf("expected project 'owner/repo', got %s", project)
	}
}

func TestDetectPlatformGitLab(t *testing.T) {
	os.Setenv("GITLAB_CI", "true")
	os.Setenv("GITLAB_TOKEN", "gl-token")
	os.Setenv("CI_PROJECT_ID", "12345")
	defer func() {
		os.Unsetenv("GITLAB_CI")
		os.Unsetenv("GITLAB_TOKEN")
		os.Unsetenv("CI_PROJECT_ID")
	}()

	platform, token, project, err := DetectPlatform()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if platform != "gitlab" {
		t.Errorf("expected platform 'gitlab', got %s", platform)
	}
	if token != "gl-token" {
		t.Errorf("expected token 'gl-token', got %s", token)
	}
	if project != "12345" {
		t.Errorf("expected project '12345', got %s", project)
	}
}

func TestDetectPlatformNoCI(t *testing.T) {
	os.Unsetenv("GITHUB_ACTIONS")
	os.Unsetenv("GITLAB_CI")

	_, _, _, err := DetectPlatform()
	if err == nil {
		t.Fatal("expected error when no CI platform detected")
	}
	if err != ErrPlatformDetectionFailed {
		t.Errorf("expected ErrPlatformDetectionFailed, got %v", err)
	}
}

func TestErrUnsupportedPlatformError(t *testing.T) {
	err := NewErrUnsupportedPlatform("custom")
	expected := "unsupported platform: custom"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

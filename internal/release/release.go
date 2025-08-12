package release

import (
	"os"
	"strings"
)

type Releaser interface {
	Release(tag, notes string) error
}

func Platform(platform, token, projectID, branch string) (Releaser, error) {
	platform = strings.ToLower(platform)
	switch platform {
	case "gitlab":
		return NewGitlabReleaser(token, projectID, branch)
	// case "github":
	// 	return NewGithubReleaser(token, projectID, branch)
	// case "gitea":
	// 	return NewGiteaReleaser(token, projectID, branch)
	default:
		return nil, NewErrUnsupportedPlatform(platform)
	}
}

func DetectPlatform() (string, error) {
	if os.Getenv("GITLAB_CI") == "true" {
		return "gitlab", nil
	}
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return "github", nil
	}
	return "", ErrPlatformDetectionFailed
}

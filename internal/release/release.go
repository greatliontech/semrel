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
	case "github":
		// split projectID into owner and repo
		parts := strings.Split(projectID, "/")
		return NewGithubReleaser(token, parts[0], parts[1], branch)
	// case "gitea":
	// 	return NewGiteaReleaser(token, projectID, branch)
	default:
		return nil, NewErrUnsupportedPlatform(platform)
	}
}

func DetectPlatform() (string, string, string, error) {
	if os.Getenv("GITLAB_CI") == "true" {
		token := os.Getenv("GITLAB_TOKEN")
		project := os.Getenv("CI_PROJECT_ID")
		return "gitlab", token, project, nil
	}
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		token := os.Getenv("GITHUB_TOKEN")
		project := os.Getenv("GITHUB_REPOSITORY")
		return "github", token, project, nil
	}
	return "", "", "", ErrPlatformDetectionFailed
}

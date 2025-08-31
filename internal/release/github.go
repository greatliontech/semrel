package release

import (
	"context"

	"github.com/google/go-github/v74/github"
)

var _ Releaser = (*githubReleaser)(nil)

type githubReleaser struct {
	client *github.Client
	owner  string
	repo   string
	branch string
}

func NewGithubReleaser(token, owner, repo, branch string) (*githubReleaser, error) {
	client := github.NewClient(nil).WithAuthToken(token)
	if branch == "" {
		repository, _, err := client.Repositories.Get(context.TODO(), owner, repo)
		if err != nil {
			return nil, err
		}
		branch = repository.GetDefaultBranch()
	}
	return &githubReleaser{
		client: client,
		owner:  owner,
		repo:   repo,
		branch: branch,
	}, nil
}

func (g *githubReleaser) Release(tag string, notes string) error {
	_, _, err := g.client.Repositories.CreateRelease(context.TODO(), g.owner, g.repo, &github.RepositoryRelease{
		TagName:         github.Ptr(tag),
		TargetCommitish: github.Ptr(g.branch),
		Name:            github.Ptr(tag),
		Body:            github.Ptr(notes),
	})
	return err
}

package release

import (
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var _ Releaser = (*gitlabReleaser)(nil)

type gitlabReleaser struct {
	client    *gitlab.Client
	projectID string
	branch    string
}

func NewGitlabReleaser(token, projectID, branch string) (*gitlabReleaser, error) {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, err
	}
	if branch == "" {
		project, _, err := client.Projects.GetProject(projectID, nil)
		if err != nil {
			return nil, err
		}
		branch = project.DefaultBranch
	}
	return &gitlabReleaser{
		client:    client,
		projectID: projectID,
		branch:    branch,
	}, nil
}

func (r *gitlabReleaser) Release(tag, notes string) error {
	_, _, err := r.client.Releases.CreateRelease(r.projectID, &gitlab.CreateReleaseOptions{
		TagName:     gitlab.Ptr(tag),
		Ref:         gitlab.Ptr(r.branch),
		Description: gitlab.Ptr(notes),
	})
	return err
}

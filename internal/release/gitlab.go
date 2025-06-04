package release

import gitlab "gitlab.com/gitlab-org/api/client-go"

type gitlabReleaser struct {
	client *gitlab.Client
}

func (r *gitlabReleaser) Release() {
	r.client.Projects.GetProject(pid any, opt *gitlab.GetProjectOptions, options ...gitlab.RequestOptionFunc)
}

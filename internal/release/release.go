package release

type Releaser interface {
	Release(tag, notes string) error
}

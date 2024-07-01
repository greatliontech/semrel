package semrel

import "github.com/Masterminds/semver/v3"

func NextVersion(current *semver.Version, commits []*Commit, cfg *Config) semver.Version {
	currentBump := BumpNone
	for _, c := range commits {
		b := c.BumpKind(cfg)
		if b.Index() > currentBump.Index() {
			currentBump = b
		}
		if currentBump.IsMajor() {
			break
		}
	}
	if currentBump.IsNone() {
		currentBump = cfg.DefaultBump
	}
	if currentBump.IsMajor() {
		return current.IncMajor()
	}
	if currentBump.IsMinor() {
		return current.IncMinor()
	}
	if currentBump.IsPatch() {
		return current.IncPatch()
	}
	return *current
}

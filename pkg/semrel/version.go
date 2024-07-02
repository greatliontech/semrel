package semrel

import "github.com/Masterminds/semver/v3"

func NextVersion(current *semver.Version, commits []*Commit, cfg *Config) semver.Version {
	currentBump := BumpNone
	for _, c := range commits {
		b := c.BumpKind(cfg)
		if b.IsGreater(currentBump) {
			currentBump = b
		}
		if currentBump == BumpMajor {
			break
		}
	}
	if currentBump == BumpNone {
		currentBump = cfg.DefaultBump()
	}
	if currentBump == BumpMajor {
		return current.IncMajor()
	}
	if currentBump == BumpMinor {
		return current.IncMinor()
	}
	if currentBump == BumpPatch {
		return current.IncPatch()
	}
	return *current
}

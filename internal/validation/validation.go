package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
)

var (
	prefixPattern     = regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)
	prereleasePattern = regexp.MustCompile(`^[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*$`)
	buildPattern      = regexp.MustCompile(`^[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*$`)
)

// ValidatePrefix validates that the prefix contains only allowed characters.
func ValidatePrefix(prefix string) error {
	if prefix == "" {
		return nil
	}
	if !prefixPattern.MatchString(prefix) {
		return fmt.Errorf("prefix must contain only alphanumeric characters, underscores, or hyphens: %q", prefix)
	}
	return nil
}

// ValidateInitialVersion validates that the initial version is a valid semver string.
func ValidateInitialVersion(version string) error {
	if version == "" {
		return nil
	}
	_, err := semver.NewVersion(version)
	if err != nil {
		return fmt.Errorf("invalid initial version %q: %w", version, err)
	}
	return nil
}

// ValidateMatchRule validates that the regex pattern is valid.
func ValidateMatchRule(pattern string) error {
	if pattern == "" {
		return fmt.Errorf("match pattern cannot be empty")
	}
	_, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid match pattern %q: %w", pattern, err)
	}
	return nil
}

// ValidatePlatform validates that the platform is one of the supported values.
func ValidatePlatform(platform string) error {
	if platform == "" {
		return nil
	}
	platform = strings.ToLower(platform)
	switch platform {
	case "github", "gitlab":
		return nil
	default:
		return fmt.Errorf("unsupported platform %q, must be 'github' or 'gitlab'", platform)
	}
}

// ValidatePrerelease validates that the prerelease string follows semver rules.
func ValidatePrerelease(prerelease string) error {
	if prerelease == "" {
		return nil
	}
	if !prereleasePattern.MatchString(prerelease) {
		return fmt.Errorf("invalid prerelease %q: must contain only alphanumeric characters, hyphens, and dots", prerelease)
	}
	return nil
}

// ValidateBuild validates that the build metadata string follows semver rules.
func ValidateBuild(build string) error {
	if build == "" {
		return nil
	}
	if !buildPattern.MatchString(build) {
		return fmt.Errorf("invalid build metadata %q: must contain only alphanumeric characters, hyphens, and dots", build)
	}
	return nil
}

// ValidateDefaultBump validates that the default bump is a valid value.
func ValidateDefaultBump(bump string) error {
	if bump == "" {
		return nil
	}
	bump = strings.ToLower(bump)
	switch bump {
	case "none", "patch", "minor", "major":
		return nil
	default:
		return fmt.Errorf("invalid default bump %q, must be one of: none, patch, minor, major", bump)
	}
}

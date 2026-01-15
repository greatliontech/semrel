package validation

import "testing"

func TestValidatePrefix(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		wantErr bool
	}{
		{"empty prefix", "", false},
		{"simple v prefix", "v", false},
		{"release prefix", "release-", false},
		{"underscore prefix", "version_", false},
		{"alphanumeric", "v123", false},
		{"invalid special char", "v@", true},
		{"invalid space", "v 1", true},
		{"invalid slash", "v/", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePrefix(tt.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePrefix(%q) error = %v, wantErr %v", tt.prefix, err, tt.wantErr)
			}
		})
	}
}

func TestValidateInitialVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{"empty version", "", false},
		{"valid version", "1.0.0", false},
		{"valid with prefix", "v1.0.0", false},
		{"valid prerelease", "1.0.0-alpha.1", false},
		{"valid with build", "1.0.0+build.123", false},
		{"invalid version", "not.a.version", true},
		{"invalid chars", "1.0.0.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInitialVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInitialVersion(%q) error = %v, wantErr %v", tt.version, err, tt.wantErr)
			}
		})
	}
}

func TestValidateMatchRule(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		wantErr bool
	}{
		{"simple pattern", "TRACK-\\d+", false},
		{"complex pattern", "^feat\\((.+)\\):", false},
		{"empty pattern", "", true},
		{"invalid regex", "[invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMatchRule(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMatchRule(%q) error = %v, wantErr %v", tt.pattern, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePlatform(t *testing.T) {
	tests := []struct {
		name     string
		platform string
		wantErr  bool
	}{
		{"empty platform", "", false},
		{"github lowercase", "github", false},
		{"github uppercase", "GITHUB", false},
		{"gitlab lowercase", "gitlab", false},
		{"gitlab mixedcase", "GitLab", false},
		{"unsupported platform", "bitbucket", true},
		{"invalid platform", "unknown", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePlatform(tt.platform)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePlatform(%q) error = %v, wantErr %v", tt.platform, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePrerelease(t *testing.T) {
	tests := []struct {
		name       string
		prerelease string
		wantErr    bool
	}{
		{"empty prerelease", "", false},
		{"simple alpha", "alpha", false},
		{"alpha with number", "alpha.1", false},
		{"rc version", "rc.1", false},
		{"beta with date", "beta-2024-01-01", false},
		{"invalid space", "alpha 1", true},
		{"invalid special char", "alpha@1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePrerelease(tt.prerelease)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePrerelease(%q) error = %v, wantErr %v", tt.prerelease, err, tt.wantErr)
			}
		})
	}
}

func TestValidateBuild(t *testing.T) {
	tests := []struct {
		name    string
		build   string
		wantErr bool
	}{
		{"empty build", "", false},
		{"simple build", "build", false},
		{"build with number", "build.123", false},
		{"sha build", "sha-abc123", false},
		{"invalid space", "build 123", true},
		{"invalid special char", "build@123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBuild(tt.build)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBuild(%q) error = %v, wantErr %v", tt.build, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDefaultBump(t *testing.T) {
	tests := []struct {
		name    string
		bump    string
		wantErr bool
	}{
		{"empty bump", "", false},
		{"none", "none", false},
		{"patch", "patch", false},
		{"minor", "minor", false},
		{"major", "major", false},
		{"uppercase PATCH", "PATCH", false},
		{"invalid value", "huge", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDefaultBump(tt.bump)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDefaultBump(%q) error = %v, wantErr %v", tt.bump, err, tt.wantErr)
			}
		})
	}
}

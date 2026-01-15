package semrel

import "testing"

func TestNewBump(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  BumpKind
		wantError bool
	}{
		{"none lowercase", "none", BumpNone, false},
		{"none uppercase", "NONE", BumpNone, false},
		{"none mixed case", "NoNe", BumpNone, false},
		{"patch lowercase", "patch", BumpPatch, false},
		{"patch uppercase", "PATCH", BumpPatch, false},
		{"minor lowercase", "minor", BumpMinor, false},
		{"minor uppercase", "MINOR", BumpMinor, false},
		{"major lowercase", "major", BumpMajor, false},
		{"major uppercase", "MAJOR", BumpMajor, false},
		{"invalid string", "invalid", BumpNone, true},
		{"empty string", "", BumpNone, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewBump(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("NewBump(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("NewBump(%q) unexpected error: %v", tt.input, err)
				return
			}
			if result != tt.expected {
				t.Errorf("NewBump(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBumpKindString(t *testing.T) {
	tests := []struct {
		bump     BumpKind
		expected string
	}{
		{BumpNone, "none"},
		{BumpPatch, "patch"},
		{BumpMinor, "minor"},
		{BumpMajor, "major"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.bump.String()
			if result != tt.expected {
				t.Errorf("BumpKind(%d).String() = %q, want %q", tt.bump, result, tt.expected)
			}
		})
	}
}

func TestBumpKindStringRoundTrip(t *testing.T) {
	bumps := []BumpKind{BumpNone, BumpPatch, BumpMinor, BumpMajor}
	for _, b := range bumps {
		str := b.String()
		result, err := NewBump(str)
		if err != nil {
			t.Errorf("round trip failed for %v: %v", b, err)
			continue
		}
		if result != b {
			t.Errorf("round trip: NewBump(%q) = %v, want %v", str, result, b)
		}
	}
}

func TestBumpKindIsGreater(t *testing.T) {
	tests := []struct {
		name     string
		b        BumpKind
		other    BumpKind
		expected bool
	}{
		{"major > minor", BumpMajor, BumpMinor, true},
		{"major > patch", BumpMajor, BumpPatch, true},
		{"major > none", BumpMajor, BumpNone, true},
		{"minor > patch", BumpMinor, BumpPatch, true},
		{"minor > none", BumpMinor, BumpNone, true},
		{"minor not > major", BumpMinor, BumpMajor, false},
		{"patch > none", BumpPatch, BumpNone, true},
		{"patch not > minor", BumpPatch, BumpMinor, false},
		{"patch not > major", BumpPatch, BumpMajor, false},
		{"none not > patch", BumpNone, BumpPatch, false},
		{"none not > minor", BumpNone, BumpMinor, false},
		{"none not > major", BumpNone, BumpMajor, false},
		{"same: major not > major", BumpMajor, BumpMajor, false},
		{"same: minor not > minor", BumpMinor, BumpMinor, false},
		{"same: patch not > patch", BumpPatch, BumpPatch, false},
		{"same: none not > none", BumpNone, BumpNone, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.b.IsGreater(tt.other)
			if result != tt.expected {
				t.Errorf("%s.IsGreater(%s) = %v, want %v", tt.b, tt.other, result, tt.expected)
			}
		})
	}
}

func TestBumpKindIsLesser(t *testing.T) {
	tests := []struct {
		name     string
		b        BumpKind
		other    BumpKind
		expected bool
	}{
		{"none < patch", BumpNone, BumpPatch, true},
		{"none < minor", BumpNone, BumpMinor, true},
		{"none < major", BumpNone, BumpMajor, true},
		{"patch < minor", BumpPatch, BumpMinor, true},
		{"patch < major", BumpPatch, BumpMajor, true},
		{"patch not < none", BumpPatch, BumpNone, false},
		{"minor < major", BumpMinor, BumpMajor, true},
		{"minor not < patch", BumpMinor, BumpPatch, false},
		{"minor not < none", BumpMinor, BumpNone, false},
		{"major not < minor", BumpMajor, BumpMinor, false},
		{"major not < patch", BumpMajor, BumpPatch, false},
		{"major not < none", BumpMajor, BumpNone, false},
		{"same: major not < major", BumpMajor, BumpMajor, false},
		{"same: minor not < minor", BumpMinor, BumpMinor, false},
		{"same: patch not < patch", BumpPatch, BumpPatch, false},
		{"same: none not < none", BumpNone, BumpNone, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.b.IsLesser(tt.other)
			if result != tt.expected {
				t.Errorf("%s.IsLesser(%s) = %v, want %v", tt.b, tt.other, result, tt.expected)
			}
		})
	}
}

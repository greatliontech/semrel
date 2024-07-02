package semrel

import (
	"fmt"
	"strings"
)

type BumpKind int

const (
	BumpNone BumpKind = iota
	BumpPatch
	BumpMinor
	BumpMajor
)

func NewBump(str string) (BumpKind, error) {
	str = strings.ToLower(str)
	switch str {
	case "none":
		return BumpNone, nil
	case "patch":
		return BumpPatch, nil
	case "minor":
		return BumpMinor, nil
	case "major":
		return BumpMajor, nil
	default:
		return BumpNone, fmt.Errorf("invalid bump kind: %s", str)
	}
}

func (b BumpKind) String() string {
	switch b {
	case BumpNone:
		return "none"
	case BumpPatch:
		return "patch"
	case BumpMinor:
		return "minor"
	case BumpMajor:
		return "major"
	default:
		panic(fmt.Sprintf("invalid bump kind: %d", b))
	}
}

func (b BumpKind) IsGreater(other BumpKind) bool {
	return b > other
}

func (b BumpKind) IsLesser(other BumpKind) bool {
	return b < other
}

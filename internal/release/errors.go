package release

import "errors"

type ErrUnsupportedPlatform struct {
	platform string
}

func NewErrUnsupportedPlatform(platform string) *ErrUnsupportedPlatform {
	return &ErrUnsupportedPlatform{platform: platform}
}

func (e *ErrUnsupportedPlatform) Error() string {
	return "unsupported platform: " + e.platform
}

var ErrPlatformDetectionFailed = errors.New("failed to detect platform")

package semrel

import "fmt"

type ConfigVersion int

const (
	ConfigVersionV0 ConfigVersion = iota
)

func NewConfigVersion(str string) (ConfigVersion, error) {
	switch str {
	case "", "v0":
		return ConfigVersionV0, nil
	default:
		return ConfigVersionV0, fmt.Errorf("invalid config version: %s", str)
	}
}

func (v ConfigVersion) String() string {
	switch v {
	case ConfigVersionV0:
		return "v0"
	default:
		panic(fmt.Sprintf("invalid config version: %d", v))
	}
}

package semrel

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrNotConventionalCommit = errors.New("not a conventional commit message")

	commitPattern   = regexp.MustCompile(`^([\w-]+)(?:\(([^\)]*)\))?(!*)\: (.*)$`)
	footerPattern   = regexp.MustCompile(`^([\w-]+): (.*)$`)
	breakingPattern = regexp.MustCompile("^BREAKING CHANGES?")
)

type Commit struct {
	Footers     map[string]string
	Type        string
	Scope       string
	Description string
	Body        string
	Attention   bool
}

func (c *Commit) BumpKind(cfg *Config) BumpKind {
	keys := make([]string, 0, len(c.Footers))

	for k := range c.Footers {
		keys = append(keys, k)
	}

	if c.Attention || FilterFooters(keys) {
		return BumpMajor
	}

	return cfg.BumpKind(c.Type)
}

func FilterFooters(keys []string) bool {
	isBreakingChange := false
	for k := range keys {
		isBreakingChange = breakingPattern.MatchString(keys[k])
	}

	return isBreakingChange
}


func ParseCommitMessage(message string) (*Commit, error) {
	lines := strings.Split(message, "\n")

	found := commitPattern.FindAllStringSubmatch(lines[0], -1)
	if len(found) < 1 {
		return nil, ErrNotConventionalCommit
	}

	attention := found[0][3] == "!"
	c := &Commit{
		Type:        strings.ToLower(found[0][1]),
		Scope:       found[0][2],
		Attention:   attention,
		Description: found[0][4],
	}

	sections := [][]string{}
	currentSection := []string{}
	for _, line := range lines[1:] {
		if line == "" {
			if len(currentSection) > 0 {
				sections = append(sections, currentSection)
				currentSection = []string{}
			}
			continue
		}
		currentSection = append(currentSection, line)
	}
	body := strings.Builder{}
	for i, section := range sections {
		if i == len(sections)-1 {
			ftr := parseFooter(section)
			if len(ftr) > 0 {
				c.Footers = ftr
				continue
			}
		}
		if i != 0 {
			body.WriteString("\n")
		}
		body.WriteString(strings.Join(section, "\n"))
	}
	body.WriteString("\n")
	c.Body = body.String()

	return c, nil
}

func parseFooter(section []string) map[string]string {
	footers := map[string]string{}
	for _, line := range section {
		found := footerPattern.FindAllStringSubmatch(line, -1)
		if len(found) < 1 {
			return nil
		}
		footers[found[0][1]] = found[0][2]
	}
	return footers
}

package semrel

import (
	"errors"
	"regexp"
	"strings"
)

type Commit struct {
	Type        string
	Scope       string
	Attention   bool
	Description string
	Body        string
	Footers     map[string]string
}

var (
	ErrNotConventionalCommit = errors.New("not a conventional commit message")

	commitPattern = regexp.MustCompile(`^([\w-]+)(?:\(([^\)]*)\))?(!*)\: (.*)$`)
	footerPattern = regexp.MustCompile(`^([\w-]+): (.*)$`)
)

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

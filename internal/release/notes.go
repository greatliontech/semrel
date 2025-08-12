package release

import (
	"regexp"
	"strings"

	"github.com/greatliontech/semrel/pkg/semrel"
)

// MatchRule holds one regex + replacement template.
type MatchRule struct {
	Match   *regexp.Regexp
	Replace string
}

// NewMatchRule creates a new MatchRule with the given regex pattern and replacement template.
func NewMatchRule(pattern, replace string) (*MatchRule, error) {
	match, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &MatchRule{
		Match:   match,
		Replace: replace,
	}, nil
}

// Apply runs the regex on s, substituting every match
// according to the Replace template (using $1, $2, … for capture groups).
func (r *MatchRule) Apply(s string) string {
	return r.Match.ReplaceAllString(s, r.Replace)
}

// Filters matches types and scopes in conventional commits, and excludes them from release notes.
type Filters struct {
	Types  []string
	Scopes []string
}

func (f *Filters) MatchType(t string) bool {
	for _, filterType := range f.Types {
		if strings.EqualFold(filterType, t) {
			return true
		}
	}
	return false
}

func (f *Filters) MatchScope(s string) bool {
	for _, filterScope := range f.Scopes {
		if strings.EqualFold(filterScope, s) {
			return true
		}
	}
	return false
}

func GenerateReleaseNotes(commits []*semrel.Commit, filters *Filters, matchRules []*MatchRule) string {
	b := strings.Builder{}
	for _, commit := range commits {
		if filters != nil &&
			(filters.MatchType(commit.Type) || filters.MatchScope(commit.Scope)) {
			continue
		}
		b.WriteString("- ")
		b.WriteString(commit.Type)
		b.WriteString("(")
		b.WriteString(commit.Scope)
		b.WriteString("): ")
		description := commit.Description
		for _, rule := range matchRules {
			description = rule.Apply(description)
		}
		b.WriteString(description)
		b.WriteString("\n")
	}
	return b.String()
}

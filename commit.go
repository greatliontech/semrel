package semrel

type Commit struct {
	Type        string
	Scope       string
	Attention   bool
	Description string
	Body        string
	Footers     map[string]string
}

func ParseCommitMessage(message string) *Commit {
	return &Commit{}
}

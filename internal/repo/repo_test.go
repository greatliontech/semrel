package repo

import (
	"math/rand"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/greatliontech/semrel"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type testCommit struct {
	msg string
	tag string
}

func testRepo(cms []testCommit) (*git.Repository, error) {
	s := memory.NewStorage()
	f := memfs.New()

	r, err := git.Init(s, f)
	if err != nil {
		return nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	for _, cm := range cms {
		file, err := f.Create(randString(10))
		if err != nil {
			return nil, err
		}
		_, err = file.Write([]byte("test content"))
		if err != nil {
			return nil, err
		}
		if err := file.Close(); err != nil {
			return nil, err
		}
		_, err = w.Add(file.Name())
		if err != nil {
			return nil, err
		}
		commit, err := w.Commit(cm.msg, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "John Doe",
				Email: "john@doe.org",
				When:  time.Now(),
			},
		})
		if err != nil {
			return nil, err
		}
		if cm.tag != "" {
			_, err = r.CreateTag(cm.tag, commit, nil)
			if err != nil {
				return nil, err
			}
		}
	}
	return r, nil
}

func TestCommits(t *testing.T) {
	commitMessages := []testCommit{
		{msg: "initial", tag: ""},
		{msg: "fix: bug", tag: ""},
		{msg: "feat: new feature", tag: ""},
	}
	r, err := testRepo(commitMessages)
	if err != nil {
		t.Fatal(err)
	}
	repo := New(r)
	commits, err := repo.Commits(plumbing.ZeroHash, plumbing.ZeroHash)
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(commits))
	}
	expected := []*semrel.Commit{
		{Type: "feat"},
		{Type: "fix"},
	}
	for i, c := range commits {
		if c.Type != expected[i].Type {
			t.Fatalf("expected %s, got %s", expected[i].Type, c.Type)
		}
	}
}

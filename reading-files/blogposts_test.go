package blogposts_test

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	blogposts "github.com/jrvldam/learn-go-with-tests/reading-files"
)

func TestNewBlogPosts(t *testing.T) {
	fs := fstest.MapFS{
		"hello world.md": {Data: []byte("Title: Post 1")},
		"hello-world.md": {Data: []byte("Title: Post 2")},
	}

	posts, err := blogposts.NewPostsFromFS(fs)
	if err != nil {
		t.Fatal(err)
	}

	got := posts[0]
	want := blogposts.Post{Title: "Post 1"}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNewBlogPostsFail(t *testing.T) {
	stub := StubFailingFs{}

	_, err := blogposts.NewPostsFromFS(stub)

	if err == nil {
		t.Errorf("expected to fail on reading a file")
	}
}

type StubFailingFs struct{}

func (s StubFailingFs) Open(name string) (fs.File, error) {
	return nil, errors.New("oh no. I always fail")
}

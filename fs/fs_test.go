package fs

// https://golang.org/pkg/testing/

import (
	// "bytes"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/afero"
)

func compare(a, b []int) bool {
	if &a == &b {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

type test struct {
	in  string
	out []int
}

func TestExist(t *testing.T) {
	appFS := afero.NewMemMapFs()
	// create test files and directories
	appFS.MkdirAll("src/a", 0755)
	afero.WriteFile(appFS, "src/a/b", []byte("file b"), 0644)
	afero.WriteFile(appFS, "src/c", []byte("file c"), 0644)
	name := "src/c"
	_, err := appFS.Stat(name)
	if os.IsNotExist(err) {
		t.Errorf("file \"%s\" does not exist.\n", name)
	}
}

func TestWalk(t *testing.T) {
	appFS := afero.NewMemMapFs()
	// create test files and directories
	appFS.MkdirAll("tmp/repo/.git", 0755)
	afero.WriteFile(appFS, "tmp/repo/.git/config", []byte(`[core]
        repositoryformatversion = 0
        filemode = true
        bare = false
        logallrefupdates = true
[remote "origin"]
        url = git@github.com:mguzelevich/repot.git
        fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
        remote = origin
        merge = refs/heads/master
`), 0644)
	afero.WriteFile(appFS, "tmp/repo/README.md", []byte("README.md"), 0644)

	r, e := WalkFs(appFS, ".")
	fmt.Printf("%v %v", r, e)
	//t.Errorf("walk", r, e)
}

package git_test

import (
	"os"
	"path/filepath"
	"testing"

	// log "github.com/sirupsen/logrus"

	"github.com/mguzelevich/repot/git"
	"github.com/mguzelevich/repot/repot_test"
)

var (
	fs           *repot_test.TestFs
	rootPath     string
	bareRepoPath string
)

func TestGit_clone(t *testing.T) {
	tfs := repot_test.NewTestFs(t)
	defer tfs.CleanUp()

	target := filepath.Join(rootPath, "tmp/test")
	client := git.NewGit(target)
	if _, err := client.Clone(bareRepoPath, target); err != nil {
		t.Error(err)
	}
	if exists := tfs.Exists(filepath.Join(target, ".git/config")); !exists {
		t.Error("file not exists")
	}
}

func TestGit_chain(t *testing.T) {
	tfs := repot_test.NewTestFs(t)
	defer tfs.CleanUp()

	target := filepath.Join(rootPath, "tmp/test1")
	client := git.NewGit(target)
	cmds := []git.GitCmd{
		git.NewGitCmd(client.Clone, []string{bareRepoPath, target}),
		git.NewGitCmd(client.Config, []string{"-l"}),
		git.NewGitCmd(client.Status, []string{}),
	}
	out, err := client.ExecChain(cmds...)
	if err != nil {
		t.Error("Error: ", err, " out: ", out)
	}
	if exists := tfs.Exists(filepath.Join(target, ".git/config")); !exists {
		t.Error("file not exists")
	}
}

func TestGit_customChain(t *testing.T) {
	tfs := repot_test.NewTestFs(t)
	defer tfs.CleanUp()

	target := filepath.Join(rootPath, "tmp/test2")
	client := git.NewGit(target)
	cmds := [][]string{
		[]string{"clone", bareRepoPath, target},
		[]string{"config", "-l"},
		[]string{"status"},
	}
	out, err := client.ExecCustomChain(cmds...)
	if err != nil {
		t.Error("Error: ", err, " out: ", out)
	}
	if exists := tfs.Exists(filepath.Join(target, ".git/config")); !exists {
		t.Error("file not exists")
	}
}

func setUp() {
	repot_test.InitLogger(nil)
	fs = repot_test.NewTestFs(nil)
	rootPath = fs.TempDir("repot~test_setup")
	bareRepoPath = filepath.Join(rootPath, "bare")
	repot_test.Untar(nil, "../testdata/bare.tar.gz", bareRepoPath)
}

func tearDown() {
	// fs.CleanUp()
}

func TestMain(m *testing.M) {
	setUp()
	defer tearDown()

	if fs.Failed() {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

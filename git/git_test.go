package git

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"

	"github.com/mguzelevich/repot/helpers"
)

var appFs = afero.NewOsFs()

func TestGit_untar(t *testing.T) {
	tmpDirPath, err := afero.TempDir(appFs, "", "repot~test")

	gzipFile, err := os.Open("test.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer gzipFile.Close()

	tmpfile, err := afero.TempFile(appFs, tmpDirPath, "repot~test")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpfile.Close()

	if err := helpers.UnzipFile(gzipFile, tmpfile); err != nil {
		t.Fatal(err)
	}

	dirPath := filepath.Join(tmpDirPath, "repository")
	if err := helpers.UntarFile(tmpfile, dirPath); err != nil {
		t.Fatal(err)
	}
}

func TestGit_clone(t *testing.T) {
	client := GitClient{"/tmp/test"}
	client.Clone("/tmp/test_repo")
}

func TestGit_chain(t *testing.T) {
	client := GitClient{"/tmp/test1"}
	out, err := client.ExecChain([]gitCmd{
		gitCmd{client.Clone, []string{"/tmp/test1"}},
		gitCmd{client.Config, []string{"-l"}},
		gitCmd{client.Status, []string{""}},
	})
	if err != nil {
		t.Error("Error ", err)
	} else {
		fmt.Fprintf(os.Stderr, "%s", out)
	}
}

package repot_test

import (
	"testing"

	"github.com/mguzelevich/repot/helpers"
)

func Untar(t *testing.T, tarGzFile string, dst string) {
	if t == nil {
		t = new(testing.T)
	}
	fs := NewTestFs(t)
	defer fs.CleanUp()

	gzipFile, err := fs.AppFs.Open(tarGzFile)
	if err != nil {
		t.Fatal(err)
	}
	defer gzipFile.Close()

	tmpDirPath := fs.TempDir("repot~test~untar")
	tmpfile := fs.TempFile(tmpDirPath, "")
	defer tmpfile.Close()

	if err := helpers.UnzipFile(gzipFile, tmpfile); err != nil {
		t.Fatal(err)
	}

	if err := helpers.UntarFile(fs.AppFs, tmpfile.Name(), dst); err != nil {
		t.Fatal(err)
	}
}

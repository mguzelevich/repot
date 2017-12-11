package helpers_test

// https://golang.org/pkg/testing/

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"

	"github.com/mguzelevich/repot/helpers"
	"github.com/mguzelevich/repot/repot_test"
)

func TestGit_untar(t *testing.T) {
	fs := repot_test.NewTestFs(t)
	defer fs.CleanUp()

	gzipFile, err := fs.AppFs.Open("../testdata/bare.tar.gz")
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

	dirPath := filepath.Join(tmpDirPath, "repository")
	if err := helpers.UntarFile(fs.AppFs, tmpfile.Name(), dirPath); err != nil {
		t.Fatal(err)
	}

	fs.Exists(filepath.Join(dirPath, "config"))
}

func Benchmark_osfs(b *testing.B) {
	fs := afero.NewOsFs()

	// f := func(fs *afero.Fs) {
	// 	gzipFile, err := afero.Open(fs, "../testdata/bare.tar.gz")
	// 	if err != nil {
	// 		b.Fatal(err)
	// 	}
	// 	defer gzipFile.Close()

	// 	tmpDirPath := afero.TempDir(fs, "", "repot~test~untar")
	// 	tmpfile := afero.TempFile(fs, tmpDirPath, "")
	// 	defer tmpfile.Close()

	// 	if err := helpers.UnzipFile(gzipFile, tmpfile); err != nil {
	// 		b.Fatal(err)
	// 	}

	// 	dirPath := filepath.Join(tmpDirPath, "repository")
	// 	if err := helpers.UntarFile(fs.AppFs, tmpfile.Name(), dirPath); err != nil {
	// 		b.Fatal(err)
	// 	}

	// 	fs.Exists(dirPath)
	// }

	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello %s", fs)
	}
}

func Benchmark_mapfs(b *testing.B) {
	fs := afero.NewMemMapFs()
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello %s", fs)
	}
}

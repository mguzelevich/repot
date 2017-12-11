package repot_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

// type TempFs struct {
// 	testFs *TestFs
// 	Files  []string
// }

// Files keeps track of files that we've used so we can clean up.
type TestFs struct {
	t     *testing.T
	Files []string
	AppFs afero.Fs
}

func NewTestFs(t *testing.T) *TestFs {
	tfs := &TestFs{t: t, AppFs: afero.NewOsFs()}
	if tfs.t == nil {
		tfs.t = new(testing.T)
	}
	return tfs
}

func (tfs *TestFs) Error(args ...interface{}) {
	tfs.t.Error(args...)
}

/*
TempDir Creates a tmp directory for us to use.
*/
func (tfs *TestFs) TempDir(dir string) string {
	name, err := afero.TempDir(tfs.AppFs, "", dir)
	if err != nil {
		tfs.Error("Failed to create the tmpDir: "+name, err)
	}

	tfs.Files = append(tfs.Files, name)
	return name
}

/*
TempFile Creates a tmp file for us to use when testing
*/
func (tfs *TestFs) TempFile(dir string, content string) afero.File {
	file, err := afero.TempFile(tfs.AppFs, dir, "file")
	if err != nil {
		tfs.Error("Failed to create the tmpFile: "+file.Name(), err)
	}

	file.WriteString(content)
	tfs.Files = append(tfs.Files, file.Name())

	return file
}

/*
FileSays returns true if the file at the path contains the expected byte array
content.
*/
func (tfs *TestFs) FileSays(path string, expected []byte) bool {
	content, err := afero.ReadFile(tfs.AppFs, path)
	if err != nil {
		tfs.Error("Failed to read file: "+path, err)
	}

	return bytes.Equal(content, expected)
}

/*
CleanUp removes all files in our test registry and calls `t.Error` if something goes
wrong.
*/
func (tfs *TestFs) CleanUp() {
	for _, path := range tfs.Files {
		if err := tfs.AppFs.RemoveAll(path); err != nil {
			tfs.Error(tfs.AppFs.Name(), err)
		}
	}

	tfs.Files = make([]string, 0)
}

/*
Exists returns true if the file exists. Calls t.Error if something goes wrong while
checking.
*/
func (tfs *TestFs) Exists(path string) bool {
	exists, err := afero.Exists(tfs.AppFs, path)
	if err != nil {
		tfs.Error("Something went wrong when checking if "+path+"exists!", err)
	}
	return exists
}

/*
DirContains returns true if the dir contains the path. Calls t.Error if
something goes wrong while checking.
*/
func (tfs *TestFs) DirContains(dir string, path string) bool {
	fullPath := filepath.Join(dir, path)
	return tfs.Exists(fullPath)
}

func (tfs *TestFs) Failed() bool {
	return tfs.t.Failed()
}

package helpers

// https://www.devdungeon.com/content/working-files-go
// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
// tar -czvf test.tar.gz -C /tmp/test_repo .

import (
	"archive/tar"
	//"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func UnzipFile(src afero.File, dst afero.File) error {
	gzipReader, err := gzip.NewReader(src)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	defer gzipReader.Close()

	// Uncompress to a writer. We'll use a file writer
	// outfileWriter, err := os.Create(dst)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer outfileWriter.Close()

	// Copy contents of gzipped file to output file
	_, err = io.Copy(dst, gzipReader)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	return nil
}

func UntarFile(fs afero.Fs, src string, dst string) error {
	srcFile, err := appFs.Open(src)
	if err != nil {
		return err
	}
	tr := tar.NewReader(srcFile)
	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer f.Close()

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}
}

package repot

import (
	"bufio"
	"os"
	// "path/filepath"
	// "time"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Manifest struct {
	Repositories []*Repository
}

func processManifestLine(line string) error {
	return nil
}

func (m *Manifest) ReadManifestFromFile(filename string) error {
	log.WithFields(log.Fields{"filename": filename}).Debug("reading config file")

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read file %q due: %v", filename, err)
	}
	log.Debug("raw config: \n%s", string(data))

	log.Debug("parsed config: %+v", m)
	return nil
}

func (m *Manifest) Add(repository string, path string, name string) error {
	log.WithFields(log.Fields{"repository": repository, "path": path, "name": name}).Debug("add to manifest")
	r := Repository{
		Repository: repository, Path: path, Name: name,
	}
	m.Repositories = append(m.Repositories, &r)
	return nil
}

func (m *Manifest) readManifestFromScanner(scanner *bufio.Scanner) error {

	for scanner.Scan() { // internally, it advances token based on sperator
		line := scanner.Text()

		r := csv.NewReader(strings.NewReader(line))
		r.Comma = ','
		r.Comment = '#'

		records, err := r.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		for _, record := range records {
			repository := record[0]
			path := record[1]
			name := record[2]
			m.Add(repository, path, name)
		}

		//fmt.Println(scanner.Bytes()) // token in bytes
	}
	return nil
}

func GetManifest(manifestFile string) (*Manifest, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	manifest := new(Manifest)
	if fi.Mode()&os.ModeNamedPipe == 0 {
		if manifestFile == "" {
			log.Debug("no pipe, no manifest file :(")
		} else {
			log.Debug("manifest file")
			fd, err := os.Open(manifestFile)
			if err != nil {
				return nil, fmt.Errorf("Open file: reading manifest error")
			}
			defer fd.Close()

			scanner = bufio.NewScanner(fd)
		}
	} else {
		if manifestFile == "" {
			log.Debug("pipe!")
		} else {
			log.Debug("pipe, manifest file skipped")
		}
	}

	if err := manifest.readManifestFromScanner(scanner); err != nil {
		return nil, fmt.Errorf("readManifestFromScanner: reading manifest error")
	}
	return manifest, nil
}

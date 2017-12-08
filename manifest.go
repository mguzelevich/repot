package repot

import (
	// "path/filepath"
	// "time"
	"encoding/csv"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/mguzelevich/repot/git"
	"github.com/mguzelevich/repot/helpers"
)

type Manifest struct {
	Repositories []*git.Repository
}

func processManifestLine(line string) error {
	return nil
}

func (m *Manifest) Add(repository string, path string, name string) error {
	log.WithFields(log.Fields{"repository": repository, "path": path, "name": name}).Debug("add to manifest")
	r := git.Repository{
		Repository: repository, Path: path, Name: name,
	}
	m.Repositories = append(m.Repositories, &r)
	return nil
}

func (m *Manifest) Parse(data string) error {
	r := csv.NewReader(strings.NewReader(data))
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
	return nil
}

func GetManifest(manifestFile string) (*Manifest, error) {
	manifestData, err := helpers.ReadData(manifestFile)
	if err != nil {
		return nil, fmt.Errorf("get manifest: ReadData error")
	}

	manifest := new(Manifest)
	if err := manifest.Parse(strings.Join(manifestData, "\n")); err != nil {
		return nil, fmt.Errorf("get manifest: Parse error")
	}
	return manifest, nil
}

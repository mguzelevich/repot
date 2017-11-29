package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/mguzelevich/repot/git"
)

var appFs = afero.NewOsFs()

func WalkFs(fs afero.Fs, rootPath string) ([]*git.Repository, error) {
	targets := []*git.Repository{}
	walker := func(walkerPath string, f os.FileInfo, err error) error {
		if f.IsDir() && filepath.Base(walkerPath) == ".git" {
			var r *git.Repository
			if walkerPath == ".git" {
				r = &git.Repository{Path: "."}
			} else {
				idx := strings.LastIndex(walkerPath, "/")
				r = &git.Repository{Path: walkerPath[:idx]}
			}

			// directory := filepath.Join(path, r.Path, r.Name)
			if config, err := git.GetGitConfig(r.Path); err != nil {
				log.WithFields(log.Fields{"repository": r}).Error("walk: get git config")
			} else {
				r.Repository = config["remote.origin.url"]
			}

			if rootPath != "." && strings.HasPrefix(r.Path, rootPath) {
				p := r.Path
				idx := strings.LastIndex(p, "/")
				r.Path = p[len(rootPath):idx]
				r.Name = p[idx+1:]
			}

			fmt.Printf("%v\n", r)
			targets = append(targets, r)
			return filepath.SkipDir
		}
		return nil
	}

	if err := afero.Walk(fs, rootPath, walker); err != nil {
		return nil, err
	}

	return targets, nil
}

func Walk(rootPath string) ([]*git.Repository, error) {
	return WalkFs(appFs, rootPath)
}

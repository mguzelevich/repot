package repos

import (
	log "github.com/sirupsen/logrus"
)

type Repos struct {

	// ManifestFile ManifestFile

	// Rules manifest file repositories descriptions
	// Rules []string

	// Targets local repositories paths
	// Targets []Repo
}

func (r *Repos) Walk() error {
	// dirs, _ := repot.Walk(r.Root)
	// r.Targets = dirs
	log.WithFields(log.Fields{"repo": r}).Debug("walk")
	return nil
}

package repot

import (
// log "github.com/sirupsen/logrus"
// "sync"
)

type Repository struct {
	// Repository git clone url
	Repository string
	// TargetPath target directory path. default: '.''
	Path string
	// TargetName target repository directory path. default: repository name from url'
	Name string
}

package git

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type Repository struct {
	// Repository git clone url
	Repository string
	// TargetPath target directory path. default: '.''
	Path string
	// TargetName target repository directory path. default: repository name from url'
	Name string
}

func (r *Repository) HashID() string {
	// uid := fmt.Sprintf("%v %s", idx, r.Repository)
	// uid, _ = repot.UUID()
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%s|%s|%s", r.Repository, r.Path, r.Name)))
	return hex.EncodeToString(hasher.Sum(nil))
}

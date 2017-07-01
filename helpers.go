package repot

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/mguzelevich/repot/git"
)

// newUUID generates a random UUID according to RFC 4122
func UUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func keys(m map[int]bool) []int {
	keys := make([]int, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys
}

func Walk(rootPath string) ([]*Repository, error) {
	targets := []*Repository{}
	walker := func(walkerPath string, f os.FileInfo, err error) error {
		if f.IsDir() && filepath.Base(walkerPath) == ".git" {
			var r *Repository
			if walkerPath == ".git" {
				r = &Repository{Path: "."}
			} else {
				idx := strings.LastIndex(walkerPath, "/")
				r = &Repository{Path: walkerPath[:idx]}
			}

			// directory := filepath.Join(path, r.Path, r.Name)
			if config, err := git.GetGitConfig(r.Path); err != nil {
				log.WithFields(log.Fields{"repository": r}).Error("walk: get git config")
			} else {
				r.Repository = config["remote.origin.url"]
				p := r.Path
				if strings.HasPrefix(p, rootPath) {
					idx := strings.LastIndex(p, "/")
					r.Path = p[len(rootPath):idx]
					r.Name = p[idx+1:]
				}

			}
			targets = append(targets, r)
			return filepath.SkipDir
		}
		return nil
	}

	if err := filepath.Walk(rootPath, walker); err != nil {
		return nil, err
	}

	return targets, nil
}

func ParseRangesString(ranges string) ([]int, error) {
	repos := map[int]bool{}

	ranges = strings.TrimSpace(ranges)
	if len(ranges) == 0 {
		return []int{}, errors.New("empty input")
	}
	segments := strings.SplitN(ranges, ",", -1)
	for _, r := range segments {
		if rr := strings.SplitN(r, "-", -1); len(rr) > 1 {
			start, start_err := strconv.Atoi(rr[0])
			finish, finish_err := strconv.Atoi(rr[1])
			if start_err != nil || finish_err != nil {
				return []int{}, errors.New("range error")
			}

			if finish <= start {
				continue
			}
			for i := start; i <= finish; i++ {
				repos[i] = true
			}
		} else {
			if val, err := strconv.Atoi(r); err != nil {
				return []int{}, err
			} else {
				repos[val] = true
			}
		}
	}
	return keys(repos), nil
}

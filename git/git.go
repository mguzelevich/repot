package git

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

type GitRepo struct {
	root string
	url  string
}

type GitClient struct {
	directory string
}

func (g *GitClient) Clone(repository string) error {
	args := []string{"git", "clone", repository, g.directory}
	out, err := ExecGitCmd(".", args)
	log.WithFields(log.Fields{"args": args, "out": out, "err": err}).Debug("git.clone")
	return err
}

func (g *GitClient) Config() error {
	args := []string{"git", "config", "-l"}

	out, err := ExecGitCmd(g.directory, args)
	log.WithFields(log.Fields{"args": args, "path": g.directory, "out": out, "err": err}).Debug("git config")
	if err != nil {
		return err
	}

	config := map[string]string{}
	for _, s := range out {
		d := strings.Split(s, `=`)
		if len(d) != 2 {
			continue
		}
		config[d[0]] = d[1]
	}
	log.WithFields(log.Fields{"config": config}).Debug("git config")
	return err
}

func Clone(repository string, directory string) ([]string, error) {
	log.WithFields(log.Fields{"repository": repository, "directory": directory}).Debug("git.clone")

	// args := []string{"journalctl", "-b", "-f"}
	args := []string{"git", "clone", repository, directory}
	//args = append(args, ".")

	out, err := ExecGitCmd(".", args)
	return out, err
}

func GetGitConfig(directory string) (map[string]string, error) {
	log.WithFields(log.Fields{"directory": directory}).Debug("get git repo config")

	// args := []string{"journalctl", "-b", "-f"}
	args := []string{"git", "config", "-l"}
	//args = append(args, ".")

	out, err := ExecGitCmd(directory, args)
	if err != nil {
		return nil, err
	}
	config := map[string]string{}
	for _, s := range out {
		d := strings.Split(s, `=`)
		if len(d) != 2 {
			continue
		}
		config[d[0]] = d[1]
	}
	log.WithFields(log.Fields{"directory": directory, "config": config}).Debug("git config")
	return config, err
}

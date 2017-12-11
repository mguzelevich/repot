package git

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

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

package git

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

type gitFunc func(args ...string) ([]string, error)

type gitCmd struct {
	f    gitFunc
	args []string
}

type gitCustomCmd struct {
	cmd  string
	args []string
}

type gitCmdOutput struct {
	args []string
	out  []string
	err  error
}

type GitClient struct {
	directory string
}

func (g *GitClient) Clone(args ...string) ([]string, error) {
	repository := args[0]
	cmdArgs := []string{"git", "clone", repository, g.directory}
	out, err := ExecGitCmd(".", cmdArgs)
	log.WithFields(log.Fields{"args": cmdArgs, "out": out, "err": err}).Debug("git clone")
	return out, err
}

func (g *GitClient) Config(args ...string) ([]string, error) {
	cmdArgs := append([]string{"git", "config"}, args...)
	out, err := ExecGitCmd(g.directory, cmdArgs)
	log.WithFields(log.Fields{"args": cmdArgs, "path": g.directory, "out": out, "err": err}).Debug("git config")
	return out, err
}

func (g *GitClient) Status(args ...string) ([]string, error) {
	cmdArgs := append([]string{"git", "status"}, args...)
	out, err := ExecGitCmd(g.directory, cmdArgs)
	log.WithFields(log.Fields{"args": cmdArgs, "path": g.directory, "out": out, "err": err}).Debug("git status")
	return out, err
}

func (g *GitClient) ExecChain(cmds []gitCmd) ([]gitCmdOutput, error) {
	for idx, cmd := range cmds {
		out, err := cmd.f(cmd.args...)
		log.WithFields(log.Fields{"cmd": cmd, "idx": idx, "out": out, "err": err}).Debug("ExecChain")
	}

	return nil, nil
}

func (g *GitClient) ExecCustomChain(cmds []gitCmd) ([]gitCmdOutput, error) {
	return nil, nil
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

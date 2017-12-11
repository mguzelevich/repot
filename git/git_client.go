package git

import (
	log "github.com/sirupsen/logrus"
	//	"github.com/spf13/afero"
)

// var appFs = afero.NewOsFs()

var changePathExcludes = map[string]bool{"clone": true}

type GitClient struct {
	directory string
}

func NewGit(directory string) *GitClient {
	return &GitClient{directory: directory}
}

func (g *GitClient) Git(cmd string, args ...string) ([]string, error) {
	rootDirectory := g.directory
	cmdArgs := append([]string{"git", cmd}, args...)
	if _, ok := changePathExcludes[cmd]; ok {
		rootDirectory = "."
	}
	out, err := ExecGitCmd(rootDirectory, cmdArgs)
	log.WithFields(log.Fields{"cmd": cmd, "args": args, "out": out, "err": err}).Debug("git")
	return out, err
}

func (g *GitClient) Clone(args ...string) ([]string, error) {
	out, err := g.Git("clone", args...)
	return out, err
}

func (g *GitClient) Config(args ...string) ([]string, error) {
	out, err := g.Git("config", args...)
	return out, err
}

func (g *GitClient) Status(args ...string) ([]string, error) {
	out, err := g.Git("status", args...)
	return out, err
}

func (g *GitClient) ExecChain(cmds ...GitCmd) ([]gitCmdOutput, error) {
	for idx, cmd := range cmds {
		out, err := cmd.f(cmd.args...)
		log.WithFields(log.Fields{"cmd": cmd, "idx": idx, "out": out, "err": err}).Debug("ExecChain")
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (g *GitClient) ExecCustomChain(cmds ...[]string) ([]gitCmdOutput, error) {
	for idx, cmd := range cmds {
		out, err := g.Git(cmd[0], cmd[1:]...)
		log.WithFields(log.Fields{"cmd": cmd, "idx": idx, "out": out, "err": err}).Debug("ExecChain")
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

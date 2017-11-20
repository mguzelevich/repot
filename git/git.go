package git

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

type GitRepo struct {
	root string
	url  string
}

func customCmdBuilder(args []string) []string {
	result := []string{"git"}
	switch args[0] {
	case "status":
		result = []string{"git", "status", "--short", "--branch"}
	default:
		result = append(result, args...)
	}
	return result
}

func customOutParser(cmd string, out []string) []string {
	result := []string{}
	switch cmd {
	case "status":
		result = out
	default:
		result = out
	}
	return result
}

// https://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
func ExecGitCmd(dir string, args []string) ([]string, error) {
	cmdPath, err := exec.LookPath(args[0])
	if err != nil {
		log.WithFields(log.Fields{"err": err, "arg": args[0]}).Error("LookPath")
	}

	cmd := exec.Cmd{
		Dir:  dir,
		Path: cmdPath,
		Args: args,
	}

	out, err := cmd.CombinedOutput()

	output := []string{}
	for _, line := range strings.Split(string(out), "\n") {
		output = append(output, line)
	}

	l := log.WithFields(log.Fields{"err": err, "out": string(out)})
	if err != nil {
		l.Error("ExecGitCmd")
		return output, fmt.Errorf("ExecGitCmd error")
	}
	l.Info("ExecGitCmd")

	return output, nil
}

func Exec(directory string, cmd []string) ([]string, error) {
	// args := []string{"journalctl", "-b", "-f"}
	args := customCmdBuilder(cmd)

	log.WithFields(log.Fields{"directory": directory, "cmd": args}).Debug("git")

	rawOut, err := ExecGitCmd(directory, args)

	out := customOutParser(args[1], rawOut)

	return out, err
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

func (g *GitRepo) Walk() error {
	// dirs, _ := repot.Walk(g.Root)
	// //	r.Targets = dirs
	log.WithFields(log.Fields{"git": g}).Debug("Walk")
	return nil
}

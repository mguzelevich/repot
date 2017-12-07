package git

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

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
	gitLogger := log.WithFields(log.Fields{"cmd": args})
	cmdPath, err := exec.LookPath(args[0])
	if err != nil {
		gitLogger.WithFields(log.Fields{"err": err}).Error("LookPath")
		return nil, fmt.Errorf("ExecGitCmd lookup error")
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

	l := gitLogger.WithFields(log.Fields{"err": err, "out": string(out)})
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

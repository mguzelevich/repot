package git

import (
// "strings"

// log "github.com/sirupsen/logrus"
)

type gitFunc func(args ...string) ([]string, error)

type GitCmd struct {
	f    gitFunc
	args []string
}

type gitCmdOutput struct {
	args []string
	out  []string
	err  error
}

func NewGitCmd(f gitFunc, args []string) GitCmd {
	return GitCmd{f: f, args: args}
}

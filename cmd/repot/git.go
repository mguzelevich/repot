package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mguzelevich/repot"
	"github.com/mguzelevich/repot/fs"
	"github.com/mguzelevich/repot/git"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git repos activity automation",
	Long:  `Git repos activity automation`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{"use": cmd.Use, "args": args}).Debug("comand called")

		rootPath := cmdArgs.Root
		if rootPath == "" {
			rootPath = "."
		}

		var fsRepos = []*git.Repository{}

		if repositories, err := fs.Walk(rootPath); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Walk")
		} else {
			fsRepos = repositories
		}

		results := newOutputs()

		supervisor := repot.NewSuperVisor(cmdArgs.Jobs)
		supervisor.ShowProgress = cmdArgs.Progress
		for _, r := range fsRepos {
			directory := filepath.Join(rootPath, r.Path, r.Name)
			repository := r.Repository

			gitFunc := func(uid string) error {
				log.WithFields(log.Fields{"uid": uid, "repository": repository, "directory": directory}).Debug("clone func")
				out, err := git.Exec(directory, args)
				results.Add(uid, out)
				return err
			}
			uid := r.HashID()
			supervisor.AddJob(uid, gitFunc)
		}
		supervisor.ExecJobs()

		for idx, r := range fsRepos {
			status := supervisor.JobState(r.HashID())
			out := results.Get(r.HashID())
			fmt.Fprintf(os.Stderr, "=== %03d === [%s] %s\n", idx+1, r.Repository, status)
			fmt.Fprintf(os.Stderr, "%s\n", strings.Join(out, "\n"))
		}
	},
}

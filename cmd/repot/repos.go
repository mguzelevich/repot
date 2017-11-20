package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mguzelevich/repot"
	"github.com/mguzelevich/repot/fs"
	"github.com/mguzelevich/repot/git"
)

// reposCmd represents the repos command
var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "Git repos activity automation",
	Long:  `Git repos activity automation`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// TODO: Work your own magic here
	// 	fmt.Println("repos called")
	// },
}

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "clone multiply repositories specified by manifest",
	Long:  `clone multiply repositories specified by manifest`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{"use": cmd.Use, "args": args}).Debug("comand called")

		rootPath := cmdArgs.Root
		if rootPath == "" {
			// t.Format(time.RFC3339Nano)
			rootPath = filepath.Join("/tmp/repot/clone", time.Now().Format("20060102_150405"))
		}

		var manifestRepos = []*git.Repository{}
		if manifest, err := repot.GetManifest(cmdArgs.Repos.ManifestFile); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("getManifest")
		} else {
			manifestRepos = manifest.Repositories
		}

		results := make(map[string][]string)

		supervisor := repot.NewSuperVisor(cmdArgs.Jobs)
		supervisor.ShowProgress = cmdArgs.Progress
		for idx, r := range manifestRepos {
			log.WithFields(log.Fields{"idx": idx, "repository": r}).Debug("get from manifest")

			directory := filepath.Join(rootPath, r.Path, r.Name)
			repository := r.Repository

			cloneFunc := func(uid string) error {
				log.WithFields(log.Fields{"uid": uid, "repository": repository, "directory": directory}).Debug("clone func")
				out, err := git.Clone(repository, directory)
				results[uid] = out
				return err
			}
			uid := r.HashID()
			supervisor.AddJob(uid, cloneFunc)
		}
		supervisor.ExecJobs()

		for idx, r := range manifestRepos {
			state := supervisor.JobState(r.HashID())
			fmt.Fprintf(os.Stderr, "=== %03d === [%s] %s\n", idx+1, r.Repository, state)
		}
	},
}

// diffCmd compare target directory and repositories specified by manifest
var diffCmd = &cobra.Command{
	Use:   "check-diff",
	Short: "compare target directory and repositories specified by manifest",
	Long:  `compare target directory and repositories specified by manifest`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{"use": cmd.Use, "args": args}).Debug("comand called")

		var manifestRepos = []*git.Repository{}
		var fsRepos = []*git.Repository{}

		if manifest, err := repot.GetManifest(cmdArgs.Repos.ManifestFile); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("getManifest")
		} else {
			manifestRepos = manifest.Repositories
		}

		if repositories, err := fs.Walk(cmdArgs.Root); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Walk")
		} else {
			fsRepos = repositories
		}

		fsMap := map[string]*git.Repository{}
		for _, r := range fsRepos {
			fsMap[r.HashID()] = r
		}

		equial := true
		for idx, r := range manifestRepos {
			rep, ok := fsMap[r.HashID()]
			equial = equial && ok
			if ok {
				if rep == nil {
					log.WithFields(log.Fields{"repository": r}).Error("WTF")
				}
				fmt.Fprintf(os.Stderr, "=== %03d === [%s] %s\n", idx+1, r.Repository, "equal")
			} else {
				fmt.Fprintf(os.Stderr, "=== %03d === [%s] %s\n", idx+1, r.Repository, "not exist")
			}
			fsMap[r.HashID()] = nil
		}

		for _, r := range fsRepos {
			if fsMap[r.HashID()] != nil {
				fmt.Fprintf(os.Stderr, "--- [%s] %s\n", r.Repository, "not in the manifest")
			}
		}

		if equial {
			log.Info("manifest == fs")
			fmt.Fprintf(os.Stderr, "\n=== SUMMARY === [%s]\n", "equial")
			os.Exit(0)
		} else {
			log.Info("manifest != fs")
			fmt.Fprintf(os.Stderr, "\n=== SUMMARY === [%s]\n", "not equial")
			os.Exit(1)
		}
	},
}

// checkCmd represents the checkCmd command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check manifest",
	Long:  `check manifest`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("%s called %s\n", cmd.Use, args)
	},
}

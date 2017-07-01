package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mguzelevich/repot"
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
			rootPath = filepath.Join("/tmp/repot/clone", time.Now().Format("20060102_150405"))
		}
		// t.Format(time.RFC3339Nano)
		// log.Printf("check called %v", rootPath)

		// cmd.Flags().Lookup("manifest").Value.String()
		if manifest, err := repot.GetManifest(cmdArgs.Repos.ManifestFile); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("getManifest")
		} else {
			supervisor := repot.NewSuperVisor(cmdArgs.Jobs)
			supervisor.ShowProgress = cmdArgs.Progress
			for idx, r := range manifest.Repositories {
				log.WithFields(log.Fields{"idx": idx, "repository": r}).Debug("get from manifest")

				directory := filepath.Join(rootPath, r.Path, r.Name)
				repository := r.Repository

				cloneFunc := func(uid string) (string, error) {
					log.WithFields(log.Fields{"uid": uid, "repository": repository, "directory": directory}).Debug("clone func")
					out, err := git.Clone(repository, directory)
					return out, err
				}
				uid := fmt.Sprintf("%v %s", idx, r.Repository)
				uid, _ = repot.UUID()
				supervisor.AddJob(uid, cloneFunc)
			}
			supervisor.ExecJobs()
		}
		//jtools.Clone(rootPath, manifest)
	},
}

// cloneCmd represents the clone command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check manifest",
	Long:  `check manifest`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("%s called %s\n", cmd.Use, args)
	},
}

// diffCmd compare target directory and repositories specified by manifest
var diffCmd = &cobra.Command{
	Use:   "check-diff",
	Short: "compare target directory and repositories specified by manifest",
	Long:  `compare target directory and repositories specified by manifest`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{"use": cmd.Use, "args": args}).Debug("comand called")

		// cmdArgs.Root
		var manifestRepos = []*repot.Repository{}
		var fsRepos = []*repot.Repository{}

		if manifest, err := repot.GetManifest(cmdArgs.Repos.ManifestFile); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("getManifest")
		} else {
			manifestRepos = manifest.Repositories
		}

		if repositories, err := repot.Walk(cmdArgs.Root); err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Walk")
		} else {
			fsRepos = repositories
		}

		manifestMap := map[string]*repot.Repository{}

		for _, r := range manifestRepos {
			manifestRepoKey := fmt.Sprintf("%v", r)
			manifestMap[manifestRepoKey] = r
		}
		equial := true
		for _, r := range fsRepos {
			localRepoKey := fmt.Sprintf("%v", r)
			if _, ok := manifestMap[localRepoKey]; !ok {
				equial = false
				log.WithFields(log.Fields{"local": localRepoKey}).Debug("diff -")
			}
			//log.WithFields(log.Fields{"idx": idx, "repository": r}).Debug("manifest")
		}

		if equial {
			log.Info("manifest == fs")
			os.Exit(0)
		} else {
			log.Info("manifest != fs")
			os.Exit(1)
		}

		// if  !=  {
		// 	log.WithFields(log.Fields{"idx": idx, "repository": r, "remote.origin.url": config["remote.origin.url"]}).Error("diff: origin url")
		// }

		//jtools.Clone(rootPath, manifest)
	},
}

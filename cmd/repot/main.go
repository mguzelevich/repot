package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/mguzelevich/repot/workerpool"
)

var cfgFile string

type Args struct {
	Debug    bool
	Progress bool
	DryRun   bool
	Jobs     int
	Root     string

	Git struct {
		RepositoriesIndexes string
	}

	Repos struct {
		ManifestFile         string
		RepositoriesIndexes  string
		WipeTargetIfConflict bool
	}
}

var (
	cmdArgs Args
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "repot",
	Short: "automation tools",
	Long: `
RepoT is a CLI tools suite for automation of development activity.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"root.debug":                 cmdArgs.Debug,
			"root.progress":              cmdArgs.Progress,
			"root.jobs":                  cmdArgs.Jobs,
			"root.root":                  cmdArgs.Root,
			"root.dry-run":               cmdArgs.DryRun,
			"git.RepositoriesIndexes":    cmdArgs.Git.RepositoriesIndexes,
			"repos.ManifestFile":         cmdArgs.Repos.ManifestFile,
			"repos.RepositoriesIndexes":  cmdArgs.Repos.RepositoriesIndexes,
			"repos.WipeTargetIfConflict": cmdArgs.Repos.WipeTargetIfConflict,
		}).Debug("root: PersistentPreRun")
		initLogger()
	},
	// PreRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("Inside rootCmd PreRun with args: %v\n", args)
	// },
	// Run: func(cmd *cobra.Command, args []string) {
	//   fmt.Printf("Inside rootCmd Run with args: %v\n", args)
	// },
	// PostRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
	// },
	// PersistentPostRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("Inside rootCmd PersistentPostRun with args: %v\n", args)
	// },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.repot.yaml)")
	RootCmd.PersistentFlags().BoolVar(&cmdArgs.Debug, "debug", false, "Enable debug mode")
	RootCmd.PersistentFlags().BoolVar(&cmdArgs.Progress, "progress", false, "Show progress")
	RootCmd.PersistentFlags().BoolVar(&cmdArgs.DryRun, "dry-run", false, "Enable dry-run mode")
	RootCmd.PersistentFlags().IntVar(&cmdArgs.Jobs, "jobs", 1, "Jobs")
	RootCmd.PersistentFlags().StringVarP(&cmdArgs.Root, "root", "r", "", "root/target directory")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.AddCommand(gitCmd)
	gitCmd.Flags().SetInterspersed(false)
	gitCmd.PersistentFlags().StringVarP(&cmdArgs.Git.RepositoriesIndexes, "filter", "f", "", "repositories to processing")

	RootCmd.AddCommand(reposCmd)
	reposCmd.PersistentFlags().StringVarP(&cmdArgs.Repos.ManifestFile, "manifest", "m", "manifest.yaml", "manifest file")
	reposCmd.PersistentFlags().StringVarP(&cmdArgs.Repos.RepositoriesIndexes, "filter", "f", "", "repositories to processing")
	reposCmd.PersistentFlags().BoolVarP(&cmdArgs.Repos.WipeTargetIfConflict, "wipe", "", false, "wipe target dir")
	reposCmd.AddCommand(cloneCmd)
	reposCmd.AddCommand(diffCmd)
}

func initLogger() {
	timestamp := time.Now().UTC().Format("20060102")

	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile(fmt.Sprintf("/tmp/repot.%s.log", timestamp), os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.ErrorLevel)
	if cmdArgs.Debug {
		log.SetLevel(log.DebugLevel)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".repot") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
	}
}

func printStatus(status string) {
	if status == "" {
		return
	}

	scale := float32(1.0)

	if width, _, err := terminal.GetSize(int(os.Stderr.Fd())); err == nil {
		if (width - 10) < len(status) {
			scale = float32(len(status)) / float32(width-10)
		}
	}

	statusOut := ""
	m := int(scale + 0.5)

	jPending := 0
	jExecuting := 0
	jFailed := 0
	jFinished := 0

	tmp := ""
	for idx, st := range status {
		if (idx+1)%m == 0 {
			statusOut += tmp
			tmp = ""
		}

		switch st {
		case ' ':
			tmp = " "
			jPending++
		case '.':
			tmp = "."
			jExecuting++
		case 'E':
			jFailed++
			tmp = "E"
		case '+':
			jFinished++
			if tmp != "E" {
				tmp = "+"
			}
		default:
			tmp = "!"
		}
	}
	percents := int(100 * float32(jFailed+jFinished) / float32(len(status)))
	fmt.Fprintf(os.Stderr, "jobs: [%s] %d/%d (%d %%)\r", statusOut, jFailed+jFinished, len(status), percents)
}

func progressLoop(wp *workerpool.WorkerPool) {
	log.Debug("status loop started")
	heartbeat := time.Tick(2 * time.Second)
	fmt.Fprintf(os.Stderr, "\n")

	if width, height, err := terminal.GetSize(int(os.Stderr.Fd())); err != nil {
		log.WithFields(log.Fields{"width": width, "height": height, "err": err}).Debug("terminal.GetSize")
	}

	for {
		select {
		case _, ok := <-wp.JobsStateChan:
			if !ok {
				return
			}
			printStatus(wp.JobsStatusString())
		case <-heartbeat:
			printStatus(wp.JobsStatusString())
		}
	}
}

func main() {
	Execute()
}

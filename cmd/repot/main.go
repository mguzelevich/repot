package main

import (
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		ManifestFile        string
		RepositoriesIndexes string
	}
}

var (
	cmdArgs Args
)

type cmdOutput struct {
	sync.Mutex
	res map[string][]string
}

func (c *cmdOutput) Add(uid string, result []string) {
	c.Lock()
	c.res[uid] = result
	c.Unlock()
}

func (c *cmdOutput) Get(uid string) []string {
	return c.res[uid]
}

func newOutputs() *cmdOutput {
	return &cmdOutput{res: make(map[string][]string)}
}

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
			"root.debug":                cmdArgs.Debug,
			"root.progress":             cmdArgs.Progress,
			"root.jobs":                 cmdArgs.Jobs,
			"root.root":                 cmdArgs.Root,
			"root.dry-run":              cmdArgs.DryRun,
			"git.RepositoriesIndexes":   cmdArgs.Git.RepositoriesIndexes,
			"repos.ManifestFile":        cmdArgs.Repos.ManifestFile,
			"repos.RepositoriesIndexes": cmdArgs.Repos.RepositoriesIndexes,
		}).Debug("root: PersistentPreRun")
		init_logger()
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

	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.repot.yaml)")
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
	reposCmd.AddCommand(cloneCmd)
	reposCmd.AddCommand(diffCmd)
}

func init_logger() {
	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetOutput(os.Stdout)
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
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	Execute()
}

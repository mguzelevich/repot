package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/mguzelevich/repot/workerpool"
)

var cfgFile string

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

	RootCmd.PersistentFlags().BoolP("debug", "", false, "Enable debug mode")
	RootCmd.PersistentFlags().BoolP("progress", "", false, "Show progress")
	RootCmd.PersistentFlags().BoolP("dry-run", "", false, "Enable dry-run mode")
	RootCmd.PersistentFlags().IntP("jobs", "", 1, "Jobs")
	RootCmd.PersistentFlags().StringP("root", "r", "", "root/target directory")

	for _, f := range []string{"debug", "progress", "dry-run", "jobs", "root"} {
		viper.BindPFlag(f, RootCmd.PersistentFlags().Lookup(f))
	}

	RootCmd.AddCommand(gitCmd)
	gitCmd.Flags().SetInterspersed(false)
	gitCmd.PersistentFlags().StringP("filter", "f", "", "repositories to processing")

	for _, f := range []string{"filter"} {
		viper.BindPFlag(f, gitCmd.PersistentFlags().Lookup(f))
	}

	RootCmd.AddCommand(reposCmd)
	reposCmd.PersistentFlags().StringP("manifest", "m", "manifest.yaml", "manifest file")
	reposCmd.PersistentFlags().StringP("filter", "f", "", "repositories to processing")
	reposCmd.AddCommand(cloneCmd)
	reposCmd.AddCommand(diffCmd)

	for _, f := range []string{"manifest", "filter"} {
		viper.BindPFlag(f, reposCmd.PersistentFlags().Lookup(f))
	}
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
	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// home, err := homedir.Dir()
		// if err != nil {
		// 	er(err)
		// }
		// viper.AddConfigPath(home)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath("$HOME")
		viper.SetConfigName(".repot")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
	}

	viper.Set("root", path.Clean(viper.GetString("root")))

	cfgJson, _ := json.Marshal(viper.AllSettings())
	fmt.Fprintf(os.Stderr, "args: %v\n", string(cfgJson))
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

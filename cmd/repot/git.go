package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// reposCmd represents the repos command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git repos activity automation",
	Long:  `Git repos activity automation`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s called %s\n", cmd.Use, args)
	},
}

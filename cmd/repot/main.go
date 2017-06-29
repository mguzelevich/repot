package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// debug := flag.Bool("debug", false, "debug mode")
	dryRun := flag.Bool("dry-run", false, "dry-run mode")

	yamlFile := flag.String("yaml", "", "yaml scenario")

	flag.Parse()

	fmt.Fprintf(os.Stderr, "unknown mode", dryRun, yamlFile)
}

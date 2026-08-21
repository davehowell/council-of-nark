package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/davehowell/council-of-nark/experiment/ecological/snapshot"
)

func main() {
	config := flag.String("config", "", "tracked ecological snapshot config")
	flag.Parse()
	if *config == "" || flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: ecological-snapshot --config <path>")
		os.Exit(2)
	}
	runner, err := snapshot.New(*config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	attempt, err := runner.Run()
	if attempt != "" {
		if rel, relErr := filepath.Rel(runner.Root, attempt); relErr == nil {
			attempt = filepath.ToSlash(rel)
		}
		fmt.Println(attempt)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

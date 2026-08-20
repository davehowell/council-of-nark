package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davehowell/council-of-nark/experiment/harness"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: council-exp <doctor|sandbox-check|freeze|plan|run|summarize|seal|verify|bundle|judge|score> [arguments]")
}
func fail(err error) { fmt.Fprintln(os.Stderr, "error:", err); os.Exit(1) }
func need(args []string, n int) error {
	if len(args) != n {
		return fmt.Errorf("expected %d positional argument(s), got %d", n, len(args))
	}
	return nil
}
func main() {
	h, err := harness.New()
	if err != nil {
		fail(err)
	}
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	command, args := os.Args[1], os.Args[2:]
	switch command {
	case "sandbox-check":
		if err := need(args, 0); err != nil {
			fail(err)
		}
		err = h.SandboxProbe()
	case "doctor":
		fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
		skip := fs.Bool("skip-model-lookup", false, "")
		if e := fs.Parse(args); e != nil {
			fail(e)
		}
		if e := need(fs.Args(), 1); e != nil {
			fail(e)
		}
		err = h.Doctor(fs.Args()[0], *skip)
	case "freeze":
		fs := flag.NewFlagSet("freeze", flag.ContinueOnError)
		config := fs.String("config", "", "tracked config path")
		if e := fs.Parse(args); e != nil {
			fail(e)
		}
		if *config == "" {
			fail(fmt.Errorf("--config is required"))
		}
		var run string
		run, err = h.FreezeRun(*config)
		if err == nil {
			fmt.Println(run)
		}
	case "plan":
		if err := need(args, 1); err != nil {
			fail(err)
		}
		err = h.Plan(args[0])
	case "run":
		fs := flag.NewFlagSet("run", flag.ContinueOnError)
		jobs := fs.Int("jobs", 0, "")
		if e := fs.Parse(args); e != nil {
			fail(e)
		}
		if e := need(fs.Args(), 1); e != nil {
			fail(e)
		}
		err = h.Run(fs.Args()[0], *jobs)
	case "summarize":
		if err := need(args, 1); err != nil {
			fail(err)
		}
		err = h.Summarize(args[0])
	case "seal":
		if err := need(args, 1); err != nil {
			fail(err)
		}
		err = h.Seal(args[0])
	case "verify":
		if err := need(args, 1); err != nil {
			fail(err)
		}
		err = h.Verify(args[0])
	case "bundle":
		if err := need(args, 1); err != nil {
			fail(err)
		}
		err = h.Bundle(args[0])
	case "judge":
		fs := flag.NewFlagSet("judge", flag.ContinueOnError)
		jobs := fs.Int("jobs", 2, "")
		config := fs.String("config", "experiment/config/judge-smoke.json", "")
		if e := fs.Parse(args); e != nil {
			fail(e)
		}
		if e := need(fs.Args(), 1); e != nil {
			fail(e)
		}
		err = h.Judge(fs.Args()[0], *config, *jobs)
	case "score":
		fs := flag.NewFlagSet("score", flag.ContinueOnError)
		label := fs.String("label", "adjudicated", "")
		if e := fs.Parse(args); e != nil {
			fail(e)
		}
		if e := need(fs.Args(), 2); e != nil {
			fail(e)
		}
		err = h.Score(fs.Args()[0], fs.Args()[1], *label)
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fail(err)
	}
}

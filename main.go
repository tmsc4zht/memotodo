package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	Version  = "unset"
	Revision = "unset"
)

func run([]string) error {
	var usageFlag bool
	var versionFlag bool

	flag.BoolVar(&usageFlag, "usage", false, "print usage")

	flag.Parse()

	if usageFlag {
		fmt.Println("open todo")
		return nil
	}

	if versionFlag {
		fmt.Printf("%s(%s)", Version, Revision)
		return nil
	}

	if len(flag.Args()) < 0 {
		return fmt.Errorf("argument required")
	}

	switch flag.Arg(0) {
	case "install":
		return (cmdInstall())
	}

	return fmt.Errorf("unknown command")
}

func cmdInstall() error {
	return fmt.Errorf("not implemented")
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	os.Exit(0)
}

package main

import (
	"flag"
	"fmt"
)

var (
	Version  = "unset"
	Revision = "unset"
)

func main() {
	var usageFlag bool
	var versionFlag bool

	flag.BoolVar(&usageFlag, "usage", false, "print usage")
	flag.BoolVar(&versionFlag, "version", false, "print version")

	flag.Parse()

	if usageFlag {
		fmt.Println("open todo")
		return
	}

	if versionFlag {
		fmt.Printf("%s(%s)", Version, Revision)
		return
	}

}

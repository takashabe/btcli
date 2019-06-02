package main

import (
	"os"
	"runtime/debug"

	"github.com/takashabe/btcli/pkg/cmd/interactive"
)

// App version
var (
	Version = "undefined"
	Sum     = "undefined"
)

func main() {
	if i, ok := debug.ReadBuildInfo(); ok {
		Version = i.Main.Version
		Sum = i.Main.Sum
	}

	cli := &interactive.CLI{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
		Version:   Version,
		Sum:       Sum,
	}
	os.Exit(cli.Run(os.Args))
}

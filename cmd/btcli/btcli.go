package main

import (
	"os"

	"github.com/takashabe/btcli/pkg/cmd/interactive"
)

// App version
var (
	Version  = "undefined"
	Revision = "undefined"
)

func main() {
	cli := &interactive.CLI{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
		Version:   Version,
		Revision:  Revision,
	}
	os.Exit(cli.Run(os.Args))
}

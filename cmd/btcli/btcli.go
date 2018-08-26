package main

import (
	"os"

	"github.com/takashabe/btcli/api/interfaces"
)

// App version
var (
	Version  = "undefined"
	Revision = "undefined"
)

func main() {
	cli := &interfaces.CLI{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
		Version:   Version,
		Revision:  Revision,
	}
	os.Exit(cli.Run(os.Args))
}

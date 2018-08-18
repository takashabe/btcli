package main

import (
	"os"

	"github.com/takashabe/btcli/api/interfaces"
)

// Version app version
var Version = "undefined"

func main() {
	cli := &interfaces.CLI{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
		Version:   Version,
	}
	os.Exit(cli.Run(os.Args))
}

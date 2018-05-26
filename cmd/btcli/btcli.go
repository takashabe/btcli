package main

import (
	"os"

	"github.com/takashabe/btcli/api/interfaces"
)

func main() {
	cli := &interfaces.CLI{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}
	os.Exit(cli.Run(os.Args))
}

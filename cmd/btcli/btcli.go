package main

import (
	"fmt"

	prompt "github.com/c-bata/go-prompt"
	"github.com/takashabe/btcli/api/interfaces"
)

func main() {
	fmt.Println("Please select table.")
	p := prompt.New(
		interfaces.Executor,
		interfaces.Completer,
	)
	p.Run()
}

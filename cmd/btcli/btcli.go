package main

import (
	"fmt"
	"os"

	prompt "github.com/c-bata/go-prompt"
)

func dummyExecutor(s string) {
	fmt.Fprintf(os.Stdout, "Write: %s", s)
}

func dummyCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "foo", Description: "bar"},
	}
	return s
}

func main() {
	fmt.Println("Please select table.")
	p := prompt.New(
		dummyExecutor,
		dummyCompleter,
	)
	p.Run()
}

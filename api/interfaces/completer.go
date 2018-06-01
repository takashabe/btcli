package interfaces

import (
	"context"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/takashabe/btcli/api/application"
)

var commands = []prompt.Suggest{
	// cbt commands
	{Text: "ls", Description: "List tables"},
	{Text: "lookup", Description: "Read from a single row"},
	{Text: "read", Description: "Read from a multi rows"},

	// btcli commands
	{Text: "exit", Description: "Exit this prompt"},
	{Text: "quit", Description: "Exit this prompt"},
}

// Completer provides completion command handler
type Completer struct {
	tableInteractor *application.TableInteractor
}

// Do provide completion to prompt
func (h *Completer) Do(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")

	return completeWithArguments(args...)
}

func completeWithArguments(args ...string) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(commands, args[0], true)
	}

	cmd := args[0]

	second := args[1]
	switch cmd {
	case "lookup", "read":
		if len(args) == 2 {
			return prompt.FilterHasPrefix(getTableSuggestions(), second, true)
		}
	}

	return []prompt.Suggest{}
}

func getTableSuggestions() []prompt.Suggest {
	tbls, err := tableInteractor.GetTables(context.Background())
	if err != nil {
		return []prompt.Suggest{}
	}

	s := make([]prompt.Suggest, 0, len(tbls))
	for _, t := range tbls {
		s = append(s, prompt.Suggest{Text: t})
	}
	return s
}

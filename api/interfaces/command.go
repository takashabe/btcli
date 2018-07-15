package interfaces

import (
	"context"

	prompt "github.com/c-bata/go-prompt"
)

// Command defines command describe and runner
type Command struct {
	Name        string
	Description string
	Usage       string
	Runner      func(context.Context, ...string)
}

var commands = []Command{
	// cbt commands
	{
		Name:        "ls",
		Description: "List tables",
		Usage:       "ls",
	},
	{
		Name:        "count",
		Description: "Count table rows",
		Usage:       "count <table>",
	},
	{
		Name:        "lookup",
		Description: "Read from a single row",
		Usage: `lookup <table> <row> [family=<column_family>] [version=<n>]
		family    Read only columns family with <columns_family>
		version   Read only latest <n> columns`,
	},
	{
		Name:        "read",
		Description: "Read from a multi rows",
		Usage: `read <table> [start=<row>] [end=<row>] [prefix=<prefix>] [family=<column_family>] [version=<n>]
		start     Start reading at this row
		end       Stop reading before this row
		prefix    Read rows with this prefix
		family    Read only columns family with <columns_family>
		version   Read only latest <n> columns`,
	},

	// btcli commands
	{
		Name:        "exit",
		Description: "Exit this prompt",
		Usage:       "Exit this prompt",
	},
	{
		Name:        "quit",
		Description: "Exit this prompt",
		Usage:       "Exit this prompt",
	},
}

func getAllSuggests() []prompt.Suggest {
	ss := make([]prompt.Suggest, 0, len(commands))
	for _, c := range commands {
		ss = append(ss, prompt.Suggest{Text: c.Name, Description: c.Description})
	}
	return ss
}

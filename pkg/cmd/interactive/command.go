package interactive

import (
	"context"

	prompt "github.com/c-bata/go-prompt"
	"github.com/takashabe/btcli/pkg/bigtable"
	"github.com/takashabe/btcli/pkg/evaluator/cbt"
)

// Command defines command describe and runner
type Command struct {
	Name        string
	Description string
	Usage       string
	Runner      func(context.Context, bigtable.Client, ...string)
}

var commands = []Command{
	{
		Name:        "help",
		Description: "help command",
		Usage:       "help [<command>]",
		Runner:      doHelp,
	},
	{
		Name:        "ls",
		Description: "List tables",
		Usage:       "ls",
		Runner:      cbt.DoLS,
	},
	{
		Name:        "count",
		Description: "Count table rows",
		Usage:       "count <table>",
		Runner:      cbt.DoCount,
	},
	{
		Name:        "lookup",
		Description: "Read from a single row",
		Usage: `lookup <table> <row> [family=<column_family>] [version=<n>]
	version        Read only latest <n> columns
	decode         Decode big-endian value
	decode-columns Decode big-endian value with columns. <column_name:<string|int|float>[,<column_name:...>]`,
		Runner: cbt.DoLookup,
	},
	{
		Name:        "read",
		Description: "Read from a multi rows",
		Usage: `read <table> [start=<row>] [end=<row>] [prefix=<prefix>] [family=<column_family>] [version=<n>]
	start          Start reading at this row
	end            Stop reading before this row
	prefix         Read rows with this prefix
	value          Read rows with has value
	family         Read only columns family with <columns_family>
	version        Read only latest <n> columns
	from           Read newer cells than this unixtime
	to             Read older cells than this unittime
	decode         Decode big-endian value
	decode-columns Decode big-endian value with columns. <column_name:<string|int|float>[,<column_name:...>]`,
		Runner: cbt.DoRead,
	},

	// btcli commands
	{
		Name:        "exit",
		Description: "Exit this prompt",
		Usage:       "Exit this prompt",
		Runner:      doExit,
	},
	{
		Name:        "quit",
		Description: "Exit this prompt",
		Usage:       "Exit this prompt",
		Runner:      doExit,
	},
}

func getAllSuggests() []prompt.Suggest {
	ss := make([]prompt.Suggest, 0, len(commands))
	for _, c := range commands {
		ss = append(ss, prompt.Suggest{Text: c.Name, Description: c.Description})
	}
	return ss
}

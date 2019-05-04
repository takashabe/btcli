package interactive

import (
	"context"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/takashabe/btcli/pkg/bigtable"
)

// Completer provides completion command handler
type Completer struct {
	client bigtable.Client
}

// Do provide completion to prompt
func (c *Completer) Do(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")

	return c.completeWithArguments(args...)
}

func (c *Completer) completeWithArguments(args ...string) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(getAllSuggests(), args[0], true)
	}

	cmd := args[0]

	second := args[1]
	switch cmd {
	case "count":
		if len(args) == 2 {
			return prompt.FilterHasPrefix(c.getTableSuggestions(), second, true)
		}
	case "lookup":
		if len(args) == 2 {
			return prompt.FilterHasPrefix(c.getTableSuggestions(), second, true)
		}

		subcommands := []prompt.Suggest{
			{Text: "version"},
		}
		if len(args) > 3 {
			distinctCommands := filterDuplicateCommands(args, subcommands)
			latestCmd := args[len(args)-1]
			return prompt.FilterHasPrefix(distinctCommands, latestCmd, true)
		}
	case "read":
		if len(args) == 2 {
			return prompt.FilterHasPrefix(c.getTableSuggestions(), second, true)
		}

		subcommands := []prompt.Suggest{
			{Text: "start"},
			{Text: "end"},
			{Text: "prefix"},
			{Text: "family"},
			{Text: "value"},
			{Text: "version"},
			{Text: "from"},
			{Text: "to"},
			{Text: "decode"},
			{Text: "decode-columns"},
		}
		if len(args) > 2 {
			distinctCommands := filterDuplicateCommands(args, subcommands)
			latestCmd := args[len(args)-1]
			return prompt.FilterHasPrefix(distinctCommands, latestCmd, true)
		}
	}

	return []prompt.Suggest{}
}

func filterDuplicateCommands(args []string, subcommands []prompt.Suggest) []prompt.Suggest {
	ret := make([]prompt.Suggest, 0)
	for _, s := range subcommands {
		exist := false
		for _, a := range args {
			if strings.HasPrefix(a, s.Text) {
				exist = true
				break
			}
		}
		if !exist {
			ret = append(ret, s)
		}
	}
	return ret
}

func (c *Completer) getTableSuggestions() []prompt.Suggest {
	tbls, err := c.client.Tables(context.Background())
	if err != nil {
		return []prompt.Suggest{}
	}

	s := make([]prompt.Suggest, 0, len(tbls))
	for _, t := range tbls {
		s = append(s, prompt.Suggest{Text: t})
	}
	return s
}

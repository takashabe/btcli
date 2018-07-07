package interfaces

import (
	"testing"

	prompt "github.com/c-bata/go-prompt"
	"github.com/stretchr/testify/assert"
)

func TestFilterDuplicateCommands(t *testing.T) {
	cases := []struct {
		args        []string
		subcommands []prompt.Suggest
		expect      []prompt.Suggest
	}{
		{
			[]string{
				"a", "b=1",
			},
			[]prompt.Suggest{
				{Text: "a"},
				{Text: "b"},
				{Text: "c"},
			},
			[]prompt.Suggest{
				{Text: "c"},
			},
		},
	}
	for _, c := range cases {
		actual := filterDuplicateCommands(c.args, c.subcommands)
		assert.Equal(t, c.expect, actual)
	}
}

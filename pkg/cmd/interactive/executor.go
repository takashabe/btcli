package interactive

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/takashabe/btcli/pkg/bigtable"
)

// Avoid to circular dependencies
var (
	doHelpFn func(context.Context, bigtable.Client, ...string)
)

func doHelp(ctx context.Context, client bigtable.Client, args ...string) {
	doHelpFn(ctx, client, args...)
}

func init() {
	doHelpFn = lazyDoHelp
}

// Executor provides exec command handler
type Executor struct {
	client  bigtable.Client
	history io.Writer
}

// Do provides execute command
func (e *Executor) Do(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}

	ctx := context.Background()
	args := strings.Split(s, " ")
	cmd := args[0]

	for _, c := range commands {
		if cmd == c.Name {
			if e.history != nil {
				fmt.Fprintln(e.history, strings.Join(args, " "))
			}

			c.Runner(ctx, e.client, args[1:]...)
			return
		}
	}
	fmt.Fprintf(e.client.ErrStream(), "Unknown command: %s\n", cmd)
}

func doExit(ctx context.Context, client bigtable.Client, args ...string) {
	fmt.Fprintln(client.OutStream(), "Bye!")
	os.Exit(0)
}

func lazyDoHelp(ctx context.Context, client bigtable.Client, args ...string) {
	if len(args) == 1 {
		usage(client.OutStream())
		return
	}
	cmd := args[1]
	for _, c := range commands {
		if c.Name == cmd {
			fmt.Fprintln(client.OutStream(), c.Usage)
			return
		}
	}
	fmt.Fprintf(client.ErrStream(), "Unknown command: %s\n", cmd)
}

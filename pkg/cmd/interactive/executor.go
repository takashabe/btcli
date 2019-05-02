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
	doHelpFn func(context.Context, *Executor, ...string)
)

func doHelp(ctx context.Context, e *Executor, args ...string) {
	doHelpFn(ctx, e, args...)
}

func init() {
	doHelpFn = lazyDoHelp
}

// Executor provides exec command handler
type Executor struct {
	outStream io.Writer
	errStream io.Writer
	history   io.Writer
}

// Do provides execute command
func (e *Executor) Do(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}

	// TODO: wip
	client, _ := bigtable.NewClient("", "",
		bigtable.WithOutStream(e.outStream),
		bigtable.WithErrStream(e.errStream),
	)

	ctx := context.Background()
	args := strings.Split(s, " ")
	cmd := args[0]

	for _, c := range commands {
		if cmd == c.Name {
			if e.history != nil {
				fmt.Fprintln(e.history, strings.Join(args, " "))
			}

			// TODO: extract args[0]
			c.Runner(ctx, client, args...)
			return
		}
	}
	fmt.Fprintf(e.errStream, "Unknown command: %s\n", cmd)
}

func doExit(ctx context.Context, e *Executor, args ...string) {
	fmt.Fprintln(e.outStream, "Bye!")
	os.Exit(0)
}

func lazyDoHelp(ctx context.Context, e *Executor, args ...string) {
	if len(args) == 1 {
		usage(e.outStream)
		return
	}
	cmd := args[1]
	for _, c := range commands {
		if c.Name == cmd {
			fmt.Fprintln(e.outStream, c.Usage)
			return
		}
	}
	fmt.Fprintf(e.errStream, "Unknown command: %s\n", cmd)
}

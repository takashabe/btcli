package interfaces

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/k0kubun/pp"
)

// Executor invoke the command
func Executor(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}

	if s == "quit" || s == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}

	ctx := context.Background()
	args := strings.Split(s, " ")
	cmd := args[0]

	// TODO: extract function per commands
	switch cmd {
	case "ls":
		tables, err := tableInteractor.GetTables(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
		pp.Println(tables)
		fmt.Fprintln(os.Stdout, pp.Sprint(tables))
	case "lookup":
		if len(args) != 3 {
			fmt.Fprintln(os.Stderr, "Invalid args: ls <table> <row>")
			return
		}
		row, err := rowsInteractor.GetRow(ctx, args[1], args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			return
		}
		fmt.Fprintln(os.Stdout, pp.Sprint(row))
	case "read":
		if len(args) != 3 {
			fmt.Fprintln(os.Stderr, "Invalid args: read <table> <prefix>")
			return
		}
		rows, err := rowsInteractor.GetRows(ctx, args[1], args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			return
		}
		fmt.Fprintln(os.Stdout, pp.Sprint(rows))
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
	}
	return
}

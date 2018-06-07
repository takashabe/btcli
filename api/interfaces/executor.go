package interfaces

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/takashabe/btcli/api/application"
	"github.com/takashabe/btcli/api/domain"
)

// Executor provides exec command handler
type Executor struct {
	tableInteractor *application.TableInteractor
	rowsInteractor  *application.RowsInteractor
}

// Do provides execute command
func (e *Executor) Do(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}

	if s == "quit" || s == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)
	}

	ctx := context.Background()
	args := strings.Split(s, " ")
	cmd := args[0]

	// TODO: extract function per commands
	switch cmd {
	case "ls":
		tables, err := e.tableInteractor.GetTables(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
		for _, tbl := range tables {
			fmt.Fprintln(os.Stdout, tbl)
		}
	case "lookup":
		if len(args) != 3 {
			fmt.Fprintln(os.Stderr, "Invalid args: ls <table> <row>")
			return
		}
		row, err := e.rowsInteractor.GetRow(ctx, args[1], args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			return
		}
		fmt.Fprintln(os.Stdout, pp.Sprint(row))
	case "read":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Invalid args: read <table> <prefix>")
			return
		}
		table := args[1]
		key := ""
		if len(args) >= 3 {
			key = args[2]
		}
		rows, err := e.rowsInteractor.GetRows(ctx, table, key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			return
		}
		printRows(rows)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
	}
	return
}

func printRows(rs []*domain.Row) {
	for _, r := range rs {
		printRow(r)
	}
}

func printRow(r *domain.Row) {
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println(r.Key)

	for _, c := range r.Columns {
		fmt.Printf("  %-40s @ %s\n", c.Qualifier, c.Version.Format("2006/01/02-15:04:05.000000"))
		fmt.Printf("    %q\n", c.Value)
	}
}

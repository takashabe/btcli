package interfaces

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/bigtable"
	"github.com/takashabe/btcli/api/application"
	"github.com/takashabe/btcli/api/domain"
)

// Executor provides exec command handler
type Executor struct {
	outStream io.Writer
	errStream io.Writer

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
		fmt.Fprintln(e.outStream, "Bye!")
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
			fmt.Fprintf(e.errStream, "%v", err)
		}
		for _, tbl := range tables {
			fmt.Fprintln(e.outStream, tbl)
		}
	case "lookup":
		if len(args) != 3 {
			fmt.Fprintln(e.errStream, "Invalid args: lookup <table> <row>")
			return
		}
		row, err := e.rowsInteractor.GetRow(ctx, args[1], args[2])
		if err != nil {
			fmt.Fprintf(e.errStream, "%v", err)
			return
		}
		e.printRow(row)
	case "read":
		if len(args) < 2 {
			fmt.Fprintln(e.errStream, "Invalid args: read <table> [args ...]")
			return
		}
		table := args[1]
		e.readWithOptions(table, args[2:]...)
	default:
		fmt.Fprintf(e.errStream, "Unknown command: %s\n", cmd)
	}
	return
}

func (e *Executor) readWithOptions(table string, args ...string) {
	parsed := make(map[string]string)
	for _, arg := range args {
		i := strings.Index(arg, "=")
		if i < 0 {
			fmt.Fprintf(os.Stderr, "Invalid args: %v\n", arg)
			return
		}
		key, val := arg[:i], arg[i+1:]
		switch key {
		default:
			fmt.Fprintf(os.Stderr, "Unknown arg: %v\n", arg)
			return
		case "prefix":
			parsed[key] = val
		}
	}

	rr, err := rowRange(parsed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invlaid range: %v\n", err)
		return
	}
	ro, err := readOption(parsed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid options: %v\n", err)
		return
	}

	ctx := context.Background()
	rows, err := e.rowsInteractor.GetRows(ctx, table, rr, ro...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}
	e.printRows(rows)
}

func rowRange(parsedArgs map[string]string) (bigtable.RowRange, error) {
	var rr bigtable.RowRange
	if prefix := parsedArgs["prefix"]; prefix != "" {
		rr = bigtable.PrefixRange(prefix)
	}

	return rr, nil
}

func readOption(parsedArgs map[string]string) ([]bigtable.ReadOption, error) {
	var opts []bigtable.ReadOption
	if count := parsedArgs["count"]; count != "" {
		n, err := strconv.ParseInt(count, 0, 64)
		if err != nil {
			return nil, err
		}
		opts = append(opts, bigtable.LimitRows(n))
	}
	if regex := parsedArgs["regex"]; regex != "" {
		opts = append(opts, bigtable.RowFilter(bigtable.RowKeyFilter(regex)))
	}

	// filter
	// TODO: decide filter option names. refs hbase-shell

	return opts, nil
}

func (e *Executor) printRows(rs []*domain.Row) {
	for _, r := range rs {
		e.printRow(r)
	}
}

func (e *Executor) printRow(r *domain.Row) {
	fmt.Fprintln(e.outStream, strings.Repeat("-", 40))
	fmt.Fprintln(e.outStream, r.Key)

	for _, c := range r.Columns {
		fmt.Fprintf(e.outStream, "  %-40s @ %s\n", c.Qualifier, c.Version.Format("2006/01/02-15:04:05.000000"))
		fmt.Fprintf(e.outStream, "    %q\n", c.Value)
	}
}

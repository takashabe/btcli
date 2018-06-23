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
		if len(args) < 3 {
			fmt.Fprintln(e.errStream, "Invalid args: lookup <table> <row>")
			return
		}
		table := args[1]
		key := args[2]
		e.lookupWithOptions(table, key, args[3:]...)
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

func (e *Executor) lookupWithOptions(table, key string, args ...string) {
	parsed := make(map[string]string)
	for _, arg := range args {
		i := strings.Index(arg, "=")
		if i < 0 {
			fmt.Fprintf(e.errStream, "Invalid args: %v\n", arg)
			return
		}
		// TODO: Improve parsing args
		k, v := arg[:i], arg[i+1:]
		switch k {
		default:
			fmt.Fprintf(e.errStream, "Unknown arg: %v\n", arg)
			return
		case "version":
			parsed[k] = v
		case "decode":
			parsed[k] = v
		case "decode_columns":
			parsed[k] = v
		}
	}

	ro, err := readOption(parsed)
	if err != nil {
		fmt.Fprintf(e.errStream, "Invalid options: %v\n", err)
		return
	}

	ctx := context.Background()
	row, err := e.rowsInteractor.GetRow(ctx, table, key, ro...)
	if err != nil {
		fmt.Fprintf(e.errStream, "%v", err)
		return
	}

	// decode options
	p := &Printer{
		outStream: e.outStream,
		errStream: e.errStream,

		decodeType:       parsed["decode"],
		decodeColumnType: decodeColumnOption(parsed),
	}
	p.printRow(row)
}

func (e *Executor) readWithOptions(table string, args ...string) {
	parsed := make(map[string]string)
	for _, arg := range args {
		i := strings.Index(arg, "=")
		if i < 0 {
			fmt.Fprintf(os.Stderr, "Invalid args: %v\n", arg)
			return
		}
		// TODO: Improve parsing args
		key, val := arg[:i], arg[i+1:]
		switch key {
		default:
			fmt.Fprintf(os.Stderr, "Unknown arg: %v\n", arg)
			return
		case "prefix":
			parsed[key] = val
		case "count":
			parsed[key] = val
		case "version":
			parsed[key] = val
		case "decode":
			parsed[key] = val
		case "decode_columns":
			parsed[key] = val
		}
	}

	rr, err := rowRange(parsed)
	if err != nil {
		fmt.Fprintf(e.errStream, "Invlaid range: %v\n", err)
		return
	}
	ro, err := readOption(parsed)
	if err != nil {
		fmt.Fprintf(e.errStream, "Invalid options: %v\n", err)
		return
	}

	ctx := context.Background()
	rows, err := e.rowsInteractor.GetRows(ctx, table, rr, ro...)
	if err != nil {
		fmt.Fprintf(e.errStream, "%v", err)
		return
	}

	// decode options
	p := &Printer{
		outStream: e.outStream,
		errStream: e.errStream,

		decodeType:       parsed["decode"],
		decodeColumnType: decodeColumnOption(parsed),
	}
	p.printRows(rows)
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
	if version := parsedArgs["version"]; version != "" {
		n, err := strconv.ParseInt(version, 0, 64)
		if err != nil {
			return nil, err
		}
		opts = append(opts, bigtable.RowFilter(bigtable.LatestNFilter(int(n))))
	}

	// TODO: Add read options. refs hbase-shell

	return opts, nil
}

func decodeColumnOption(parsedArgs map[string]string) map[string]string {
	arg := parsedArgs["decode_columns"]
	if len(arg) == 0 {
		return map[string]string{}
	}

	ds := strings.Split(arg, ",")
	ret := map[string]string{}
	for _, d := range ds {
		ct := strings.SplitN(d, ":", 2)
		if len(ct) != 2 {
			continue
		}
		ret[ct[0]] = ct[1]
	}
	return ret
}

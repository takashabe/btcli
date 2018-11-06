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

	tableInteractor *application.TableInteractor
	rowsInteractor  *application.RowsInteractor
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

			// TODO: extract args[0]
			c.Runner(ctx, e, args...)
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

func doLS(ctx context.Context, e *Executor, args ...string) {
	tables, err := e.tableInteractor.GetTables(ctx)
	if err != nil {
		fmt.Fprintf(e.errStream, "%v", err)
		return
	}
	for _, tbl := range tables {
		fmt.Fprintln(e.outStream, tbl)
	}
}

func doCount(ctx context.Context, e *Executor, args ...string) {
	if len(args) < 2 {
		fmt.Fprintln(e.errStream, "Invalid args: count <table>")
		return
	}
	table := args[1]
	cnt, err := e.rowsInteractor.GetRowCount(ctx, table)
	if err != nil {
		fmt.Fprintf(e.errStream, "%v", err)
		return
	}
	fmt.Fprintln(e.outStream, cnt)
}

func doLookup(ctx context.Context, e *Executor, args ...string) {
	if len(args) < 3 {
		fmt.Fprintln(e.errStream, "Invalid args: lookup <table> <row>")
		return
	}
	table := args[1]
	key := args[2]
	e.lookupWithOptions(table, key, args[3:]...)
}

func doRead(ctx context.Context, e *Executor, args ...string) {
	if len(args) < 2 {
		fmt.Fprintln(e.errStream, "Invalid args: read <table> [args ...]")
		return
	}
	table := args[1]
	e.readWithOptions(table, args[2:]...)
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
		case "decode", "decode_columns":
			parsed[k] = v
		case "version":
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
		case "decode", "decode_columns":
			parsed[key] = val
		case "count", "start", "end", "prefix", "version", "family":
			parsed[key] = val
		}
	}

	if (parsed["start"] != "" || parsed["end"] != "") && parsed["prefix"] != "" {
		fmt.Fprintf(e.errStream, `"start"/"end" may not be mixed with "prefix"`)
		return
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
		fmt.Fprintf(e.errStream, "%v\n", err)
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
	if start, end := parsedArgs["start"], parsedArgs["end"]; end != "" {
		rr = bigtable.NewRange(start, end)
	} else if start != "" {
		rr = bigtable.InfiniteRange(start)
	}
	if prefix := parsedArgs["prefix"]; prefix != "" {
		rr = bigtable.PrefixRange(prefix)
	}

	return rr, nil
}

func readOption(parsedArgs map[string]string) ([]bigtable.ReadOption, error) {
	var (
		opts []bigtable.ReadOption
		fils []bigtable.Filter
	)

	// filters
	if regex := parsedArgs["regex"]; regex != "" {
		// opts = append(opts, bigtable.RowFilter(bigtable.RowKeyFilter(regex)))
		fils = append(fils, bigtable.RowKeyFilter(regex))
	}
	if family := parsedArgs["family"]; family != "" {
		// opts = append(opts, bigtable.RowFilter(bigtable.FamilyFilter(fmt.Sprintf("^%s$", family))))
		fils = append(fils, bigtable.FamilyFilter(fmt.Sprintf("^%s$", family)))
	}
	if version := parsedArgs["version"]; version != "" {
		n, err := strconv.ParseInt(version, 0, 64)
		if err != nil {
			return nil, err
		}
		// opts = append(opts, bigtable.RowFilter(bigtable.LatestNFilter(int(n))))
		fils = append(fils, bigtable.LatestNFilter(int(n)))
	}

	if len(fils) == 1 {
		opts = append(opts, bigtable.RowFilter(fils[0]))
	} else if len(fils) > 1 {
		opts = append(opts, bigtable.RowFilter(bigtable.ChainFilters(fils...)))
	}

	// isolated readOption
	if count := parsedArgs["count"]; count != "" {
		n, err := strconv.ParseInt(count, 0, 64)
		if err != nil {
			return nil, err
		}
		opts = append(opts, bigtable.LimitRows(n))
	}
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

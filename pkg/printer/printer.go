package printer

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/takashabe/btcli/pkg/domain"
)

// decode to specific type.
const (
	DecodeTypeString = "string"
	DecodeTypeInt    = "int"
	DecodeTypeFloat  = "float"
)

// Printer print the bigtable items to stream
type Printer struct {
	OutStream        io.Writer
	DecodeType       string
	DecodeColumnType map[string]string
}

// PrintRows prints the list of values.
func (w *Printer) PrintRows(rs []*domain.Row) {
	for _, r := range rs {
		w.PrintRow(r)
	}
}

// PrintRow prints the value.
func (w *Printer) PrintRow(r *domain.Row) {
	fmt.Fprintln(w.OutStream, strings.Repeat("-", 40))
	fmt.Fprintln(w.OutStream, r.Key)

	for _, c := range r.Columns {
		fmt.Fprintf(w.OutStream, "  %-40s @ %s\n", c.Qualifier, c.Version.Format("2006/01/02-15:04:05.000000"))
		w.printValue(c.Qualifier, c.Value)
	}
}

func (w *Printer) printValue(q string, v []byte) {
	// extract columnName in a qualifier
	// qualifier format: "columnFamily:columnName"
	q = q[strings.Index(q, ":")+1:]

	// retrieve decode each columns
	// decodeColumns format "column1:type1,column2:type2,..."
	for column, decode := range w.DecodeColumnType {
		if q == column {
			w.doPrint(decode, v)
			return
		}
	}

	// invoke print with a general DecodeType
	w.doPrint(w.DecodeType, v)
}

func (w *Printer) doPrint(decode string, v []byte) {
	if len(v) != 8 {
		fmt.Fprintf(w.OutStream, "    %q\n", v)
		return
	}

	switch decode {
	case DecodeTypeInt:
		fmt.Fprintf(w.OutStream, "    %d\n", w.byte2Int(v))
	case DecodeTypeFloat:
		fmt.Fprintf(w.OutStream, "    %f\n", w.byte2Float(v))
	case DecodeTypeString:
	default:
		fmt.Fprintf(w.OutStream, "    %q\n", v)
	}
}

func (*Printer) byte2Int(b []byte) int64 {
	return (int64)(binary.BigEndian.Uint64(b))
}

func (*Printer) byte2Float(b []byte) float64 {
	bits := binary.BigEndian.Uint64(b)
	return math.Float64frombits(bits)
}

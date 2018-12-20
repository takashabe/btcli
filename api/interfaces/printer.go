package interfaces

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/takashabe/btcli/api/domain"
)

const (
	decodeTypeString = "string"
	decodeTypeInt    = "int"
	decodeTypeFloat  = "float"
)

// Printer print the bigtable items to stream
type Printer struct {
	outStream io.Writer
	errStream io.Writer

	decodeType       string
	decodeColumnType map[string]string
}

func (w *Printer) printRows(rs []*domain.Row) {
	for _, r := range rs {
		w.printRow(r)
	}
}

func (w *Printer) printRow(r *domain.Row) {
	fmt.Fprintln(w.outStream, strings.Repeat("-", 40))
	fmt.Fprintln(w.outStream, r.Key)

	for _, c := range r.Columns {
		fmt.Fprintf(w.outStream, "  %-40s @ %s\n", c.Qualifier, c.Version.Format("2006/01/02-15:04:05.000000"))
		w.printValue(c.Qualifier, c.Value)
	}
}

func (w *Printer) printValue(q string, v []byte) {
	// extract columnName in a qualifier
	// qualifier format: "columnFamily:columnName"
	q = q[strings.Index(q, ":")+1:]

	// retrieve decode each columns
	// decodeColumns format "column1:type1,column2:type2,..."
	for column, decode := range w.decodeColumnType {
		if q == column {
			w.doPrint(decode, v)
			return
		}
	}

	// invoke print with a general decodeType
	w.doPrint(w.decodeType, v)
}

func (w *Printer) doPrint(decode string, v []byte) {
	if len(v) != 8 {
		fmt.Fprintf(w.outStream, "    %q\n", v)
		return
	}

	switch decode {
	case decodeTypeString:
		fmt.Fprintf(w.outStream, "    %q\n", v)
	case decodeTypeInt:
		fmt.Fprintf(w.outStream, "    %d\n", w.byte2Int(v))
	case decodeTypeFloat:
		fmt.Fprintf(w.outStream, "    %f\n", w.byte2Float(v))
	default:
		w.doGuessPrint(v)
	}
}

func (w *Printer) doGuessPrint(v []byte) {
	if len(v) != 8 {
		fmt.Fprintf(w.outStream, "    %q\n", v)
		return
	}

	// guess: float decides by high 2-bit flag
	// https://en.wikipedia.org/wiki/Double-precision_floating-point_format
	switch v[0] << 1 >> 7 & 1 {
	case 1:
		fmt.Fprintf(w.outStream, "    %f\n", w.byte2Float(v))
	case 0:
		fmt.Fprintf(w.outStream, "    %d\n", w.byte2Int(v))
	}
}

func (*Printer) byte2Int(b []byte) int64 {
	return (int64)(binary.BigEndian.Uint64(b))
}

func (*Printer) byte2Float(b []byte) float64 {
	bits := binary.BigEndian.Uint64(b)
	return math.Float64frombits(bits)
}

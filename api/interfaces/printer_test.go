package interfaces

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takashabe/btcli/api/domain"
)

func TestPrintRows(t *testing.T) {
	cases := []struct {
		input  *domain.Row
		expect string
	}{
		{
			&domain.Row{
				Key: "a",
				Columns: []*domain.Column{
					{
						Family:    "d",
						Qualifier: "d:row",
						Value:     []byte("a1"),
					},
				},
			},
			"----------------------------------------\na\n  d:row                                    @ 0001/01/01-00:00:00.000000\n    \"a1\"\n",
		},
	}
	for _, c := range cases {
		var buf bytes.Buffer
		printer := &Printer{
			outStream: &buf,
			errStream: &buf,
		}

		printer.printRow(c.input)
		assert.Equal(t, c.expect, buf.String())
	}
}

func TestPrintValue(t *testing.T) {
	cases := []struct {
		printer   *Printer
		qualifier string
		value     []byte
		expect    string
	}{
		{
			// decode string
			&Printer{
				decodeType: "string",
				decodeColumnType: map[string]string{
					"r":  "int",
					"ro": "float",
				},
			},
			"d:row",
			[]byte("a"),
			`"a"`,
		},
		{
			// decode float
			&Printer{
				decodeType: "string",
				decodeColumnType: map[string]string{
					"r":  "int",
					"ro": "float",
				},
			},
			"d:ro",
			[]byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 2.0
			"2.000000",
		},
		{
			// decode int
			&Printer{
				decodeType: "string",
				decodeColumnType: map[string]string{
					"r":  "int",
					"ro": "float",
				},
			},
			"d:r",
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, // 1
			"1",
		},
		{
			// decode string
			&Printer{},
			"d:row",
			[]byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 2.0
			"\"@\\x00\\x00\\x00\\x00\\x00\\x00\\x00\"",
		},
	}
	for _, c := range cases {
		var buf bytes.Buffer
		c.printer.outStream = &buf
		c.printer.errStream = &buf

		c.printer.printValue(c.qualifier, c.value)
		assert.Equal(t, c.expect, strings.TrimSpace(buf.String()))
	}
}

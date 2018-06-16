package interfaces

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"cloud.google.com/go/bigtable"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/takashabe/btcli/api/application"
	"github.com/takashabe/btcli/api/domain"
	"github.com/takashabe/btcli/api/domain/repository"
)

func TestRowRange(t *testing.T) {
	cases := []struct {
		input  map[string]string
		expect bigtable.RowRange
	}{
		{
			map[string]string{
				"prefix": "1",
			},
			bigtable.NewRange("1", "2"),
		},
		{
			map[string]string{
				"start": "1",
				"end":   "2",
			},
			bigtable.NewRange("", ""),
		},
	}
	for _, c := range cases {
		actual, err := rowRange(c.input)
		assert.NoError(t, err)
		assert.Equal(t, c.expect, actual)
	}
}

func TestReadOption(t *testing.T) {
	cases := []struct {
		input   map[string]string
		expects []bigtable.ReadOption
	}{
		{
			map[string]string{
				"count": "1",
			},
			[]bigtable.ReadOption{
				bigtable.LimitRows(0),
			},
		},
		{
			map[string]string{
				"count": "1",
				"regex": "1",
			},
			[]bigtable.ReadOption{
				bigtable.LimitRows(0),
				bigtable.RowFilter(bigtable.RowKeyFilter("")),
			},
		},
	}
	for _, c := range cases {
		actual, err := readOption(c.input)
		assert.NoError(t, err)

		for _, e := range c.expects {
			contain := false
			expectType := reflect.TypeOf(e)
			for _, a := range actual {
				if expectType == reflect.TypeOf(a) {
					contain = true
				}
			}
			if !contain {
				assert.Fail(t, fmt.Sprintf("Expect contan type '%v'", expectType))
			}
		}
	}
}

func TestDoExecutor(t *testing.T) {
	cases := []struct {
		input   string
		expect  string
		prepare func(*repository.MockBigtable)
	}{
		{
			"ls",
			"a\nb\n",
			func(mock *repository.MockBigtable) {
				mock.EXPECT().Tables(gomock.Any()).Return([]string{"a", "b"}, nil).Times(1)
			},
		},
		{
			"lookup table a",
			"----------------------------------------\na\n  d:row                                    @ 0001/01/01-00:00:00.000000\n    \"a1\"\n",
			func(mock *repository.MockBigtable) {
				mock.EXPECT().Get(gomock.Any(), "table", "a").Return(
					&domain.Bigtable{
						Table: "table",
						Rows: []*domain.Row{
							&domain.Row{
								Key: "a",
								Columns: []*domain.Column{
									&domain.Column{
										Family:    "d",
										Qualifier: "d:row",
										Value:     []byte("a1"),
									},
								},
							},
						},
					}, nil).Times(1)
			},
		},
		{
			"read table prefix=a",
			"----------------------------------------\na\n  d:row                                    @ 0001/01/01-00:00:00.000000\n    \"a1\"\n",
			func(mock *repository.MockBigtable) {
				mock.EXPECT().GetRows(gomock.Any(), "table", bigtable.PrefixRange("a")).Return(
					&domain.Bigtable{
						Table: "table",
						Rows: []*domain.Row{
							&domain.Row{
								Key: "a",
								Columns: []*domain.Column{
									&domain.Column{
										Family:    "d",
										Qualifier: "d:row",
										Value:     []byte("a1"),
									},
								},
							},
						},
					}, nil).Times(1)
			},
		},
	}
	for _, c := range cases {
		ctrl := gomock.NewController(t)
		mockBtRepo := repository.NewMockBigtable(ctrl)
		defer ctrl.Finish()

		c.prepare(mockBtRepo)

		var buf bytes.Buffer
		// TODO: debug
		// var r io.Reader = &buf
		// r = io.TeeReader(r, os.Stdout)
		executor := Executor{
			outStream:       &buf,
			errStream:       &buf,
			tableInteractor: application.NewTableInteractor(mockBtRepo),
			rowsInteractor:  application.NewRowsInteractor(mockBtRepo),
		}

		executor.Do(c.input)
		assert.Equal(t, c.expect, buf.String())
	}
}

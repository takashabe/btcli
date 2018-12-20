package interfaces

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

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
			bigtable.NewRange("1", "2"),
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
				bigtable.LimitRows(1),
			},
		},
		{
			map[string]string{
				"count": "1",
				"regex": "a",
			},
			[]bigtable.ReadOption{
				bigtable.RowFilter(bigtable.RowKeyFilter("a")),
				bigtable.LimitRows(1),
			},
		},
		{
			map[string]string{
				"family":  "d",
				"version": "1",
				"from":    "1545000981",
				"to":      "1545100981",
				"value":   "a$",
			},
			[]bigtable.ReadOption{
				bigtable.RowFilter(bigtable.ChainFilters(
					bigtable.FamilyFilter("^d$"),
					bigtable.LatestNFilter(1),
					bigtable.TimestampRangeFilter(time.Unix(1545000981, 0), time.Unix(1545100981, 0)),
					bigtable.ValueFilter("a$"),
				)),
			},
		},
	}
	for _, c := range cases {
		actual, err := readOption(c.input)
		assert.NoError(t, err)
		assert.Equal(t, c.expects, actual)
	}
}

func TestDoReadRowExecutor(t *testing.T) {
	tm, _ := time.Parse("2006-01-02 15:04:05", "2018-01-01 00:00:00")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		env     map[string]string
		input   string
		expect  string
		prepare func(*repository.MockBigtable)
	}{
		{
			map[string]string{},
			"ls",
			"a\nb\n",
			func(mock *repository.MockBigtable) {
				mock.EXPECT().Tables(gomock.Any()).Return([]string{"a", "b"}, nil).Times(1)
			},
		},
		{
			map[string]string{},
			"lookup table a version=1 decode=int decode_columns=row:string,404:float",
			"----------------------------------------\na\n  d:row                                    @ 2018/01/01-00:00:00.000000\n    \"a1\"\n",
			func(mock *repository.MockBigtable) {
				mock.EXPECT().Get(
					gomock.Any(),
					"table",
					"a",
					bigtable.RowFilter(bigtable.LatestNFilter(1)),
				).Return(
					&domain.Bigtable{
						Table: "table",
						Rows: []*domain.Row{
							{
								Key: "a",
								Columns: []*domain.Column{
									{
										Family:    "d",
										Qualifier: "d:row",
										Value:     []byte("a1"),
										Version:   tm,
									},
								},
							},
						},
					}, nil).Times(1)
			},
		},
		{
			map[string]string{},
			"read table prefix=a version=1 decode=int decode_columns=row:string,404:float",
			"----------------------------------------\na\n  d:row                                    @ 2018/01/01-00:00:00.000000\n    \"a1\"\n",
			func(mock *repository.MockBigtable) {
				mock.EXPECT().GetRows(
					gomock.Any(),
					"table",
					bigtable.PrefixRange("a"),
					bigtable.RowFilter(bigtable.LatestNFilter(1)),
				).Return(
					&domain.Bigtable{
						Table: "table",
						Rows: []*domain.Row{
							{
								Key: "a",
								Columns: []*domain.Column{
									{
										Family:    "d",
										Qualifier: "d:row",
										Value:     []byte("a1"),
										Version:   tm,
									},
								},
							},
						},
					}, nil).Times(1)
			},
		},
		{
			map[string]string{},
			"read table version=1 family=d",
			"----------------------------------------\na\n  d:row                                    @ 2018/01/01-00:00:00.000000\n    \"a1\"\n",
			func(mock *repository.MockBigtable) {
				mock.EXPECT().GetRows(
					gomock.Any(),
					"table",
					bigtable.RowRange{},
					chainFilters(
						bigtable.FamilyFilter("^d$"),
						bigtable.LatestNFilter(1),
					),
				).Return(
					&domain.Bigtable{
						Table: "table",
						Rows: []*domain.Row{
							{
								Key: "a",
								Columns: []*domain.Column{
									{
										Family:    "d",
										Qualifier: "d:row",
										Value:     []byte("a1"),
										Version:   tm,
									},
								},
							},
						},
					}, nil).Times(1)
			},
		},
		{
			map[string]string{
				"BTCLI_DECODE_TYPE": "int",
			},
			"read table version=1 family=d",
			"----------------------------------------\na\n  d:row                                    @ 2018/01/01-00:00:00.000000\n    1\n",
			func(mock *repository.MockBigtable) {
				mock.EXPECT().GetRows(
					gomock.Any(),
					"table",
					bigtable.RowRange{},
					chainFilters(
						bigtable.FamilyFilter("^d$"),
						bigtable.LatestNFilter(1),
					),
				).Return(
					&domain.Bigtable{
						Table: "table",
						Rows: []*domain.Row{
							{
								Key: "a",
								Columns: []*domain.Column{
									{
										Family:    "d",
										Qualifier: "d:row",
										Value:     []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
										Version:   tm,
									},
								},
							},
						},
					}, nil).Times(1)
			},
		},
	}
	for _, c := range cases {
		mockBtRepo := repository.NewMockBigtable(ctrl)
		c.prepare(mockBtRepo)

		var buf bytes.Buffer
		// TODO: debug
		var r io.Reader = &buf
		r = io.TeeReader(r, os.Stdout)
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

func TestDoCountExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		input   string
		expect  string
		prepare func(*repository.MockBigtable)
	}{
		{
			"count table",
			"1\n",
			func(mock *repository.MockBigtable) {
				mock.EXPECT().Count(gomock.Any(), "table").Return(1, nil)
			},
		},
	}
	for _, c := range cases {
		mockBtRepo := repository.NewMockBigtable(ctrl)
		c.prepare(mockBtRepo)

		var buf bytes.Buffer
		// TODO: debug
		// var r io.Reader = &buf
		// r = io.TeeReader(r, os.Stdout)
		executor := Executor{
			outStream:      &buf,
			errStream:      &buf,
			rowsInteractor: application.NewRowsInteractor(mockBtRepo),
		}

		executor.Do(c.input)
		assert.Equal(t, c.expect, buf.String())
	}
}

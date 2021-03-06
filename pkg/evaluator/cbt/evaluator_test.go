package cbt

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	bt "github.com/takashabe/btcli/pkg/bigtable"
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

func TestDoRead(t *testing.T) {
	tm, _ := time.Parse("2006-01-02 15:04:05", "2018-01-01 00:00:00")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		env     map[string]string
		input   []string
		expect  string
		prepare func(*bt.MockClient)
	}{
		{
			map[string]string{},
			[]string{
				"table", "prefix=a", "version=1", "decode=int", "decode_columns=row:string,404:float",
			},
			"----------------------------------------\na\n  d:row                                    @ 2018/01/01-00:00:00.000000\n    \"a1\"\n",
			func(mock *bt.MockClient) {
				mock.EXPECT().GetRows(
					gomock.Any(),
					"table",
					bigtable.PrefixRange("a"),
					bigtable.RowFilter(bigtable.LatestNFilter(1)),
				).Return(
					&bt.Bigtable{
						Table: "table",
						Rows: []*bt.Row{
							{
								Key: "a",
								Columns: []*bt.Column{
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
			[]string{
				"table", "version=1", "family=d",
			},
			"----------------------------------------\na\n  d:row                                    @ 2018/01/01-00:00:00.000000\n    1\n",
			func(mock *bt.MockClient) {
				mock.EXPECT().GetRows(
					gomock.Any(),
					"table",
					bigtable.RowRange{},
					filtersToReadOption(
						bigtable.FamilyFilter("^d$"),
						bigtable.LatestNFilter(1),
					),
				).Return(
					&bt.Bigtable{
						Table: "table",
						Rows: []*bt.Row{
							{
								Key: "a",
								Columns: []*bt.Column{
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
		{
			map[string]string{
				"BTCLI_DECODE_TYPE": "string",
			},
			[]string{
				"table", "version=1", "family=d", "decode=int",
			},
			"----------------------------------------\na\n  d:row                                    @ 2018/01/01-00:00:00.000000\n    1\n",
			func(mock *bt.MockClient) {
				mock.EXPECT().GetRows(
					gomock.Any(),
					"table",
					bigtable.RowRange{},
					filtersToReadOption(
						bigtable.FamilyFilter("^d$"),
						bigtable.LatestNFilter(1),
					),
				).Return(
					&bt.Bigtable{
						Table: "table",
						Rows: []*bt.Row{
							{
								Key: "a",
								Columns: []*bt.Column{
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
		mockClient := bt.NewMockClient(ctrl)
		c.prepare(mockClient)

		for k, v := range c.env {
			os.Setenv(k, v)
			defer os.Setenv(k, "")
		}

		var buf bytes.Buffer
		mockClient.EXPECT().OutStream().Return(&buf).AnyTimes()
		mockClient.EXPECT().ErrStream().Return(&buf).AnyTimes()

		DoRead(context.Background(), mockClient, c.input...)
		assert.Equal(t, c.expect, buf.String())
	}
}

func TestDoCountExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		input   []string
		expect  string
		prepare func(*bt.MockClient)
	}{
		{
			[]string{"table"},
			"1\n",
			func(mock *bt.MockClient) {
				mock.EXPECT().Count(gomock.Any(), "table").Return(1, nil)
			},
		},
	}
	for _, c := range cases {
		mockClient := bt.NewMockClient(ctrl)
		c.prepare(mockClient)

		var buf bytes.Buffer
		mockClient.EXPECT().OutStream().Return(&buf).AnyTimes()
		mockClient.EXPECT().ErrStream().Return(&buf).AnyTimes()

		DoCount(context.Background(), mockClient, c.input...)
		assert.Equal(t, c.expect, buf.String())
	}
}

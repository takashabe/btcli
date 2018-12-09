package bigtable

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/stretchr/testify/assert"
	"github.com/takashabe/btcli/api/domain"
)

func TestGet(t *testing.T) {
	loadFixture(t, "testdata/users.yaml")
	loadFixture(t, "testdata/articles.yaml")
	now := time.Now()

	cases := []struct {
		table  string
		key    string
		expect *domain.Row
	}{
		{
			"users",
			"1",
			&domain.Row{
				Key: "1",
				Columns: []*domain.Column{
					{
						Family:    "d",
						Qualifier: "d:row",
						Value:     []byte("madoka"),
						Version:   now,
					},
				},
			},
		},
		{
			"articles",
			"1##1",
			&domain.Row{
				Key: "1##1",
				Columns: []*domain.Column{
					{
						Family:    "d",
						Qualifier: "d:content",
						Value:     []byte("madoka_content"),
						Version:   now,
					},
					{
						Family:    "d",
						Qualifier: "d:title",
						Value:     []byte("madoka_title"),
						Version:   now,
					},
				},
			},
		},
	}
	for _, c := range cases {
		r, err := NewBigtableRepository("test-project", "test-instance")
		assert.NoError(t, err)

		bt, err := r.Get(context.Background(), c.table, c.key)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(bt.Rows))
		actual := bt.Rows[0]
		// NOTE: hack to version timestamp
		for _, co := range actual.Columns {
			co.Version = now
		}
		assert.Equal(t, c.expect, actual)
	}
}

func TestGetRows(t *testing.T) {
	loadFixture(t, "testdata/users.yaml")
	loadFixture(t, "testdata/articles.yaml")
	tm, _ := time.Parse("2006-01-02 15:04:05", "2018-01-01 00:00:00")
	tm = tm.Local()

	cases := []struct {
		table  string
		rr     bigtable.RowRange
		opts   []bigtable.ReadOption
		expect []*domain.Row
	}{
		{
			"users",
			bigtable.PrefixRange("1"),
			[]bigtable.ReadOption{},
			[]*domain.Row{
				{
					Key: "1",
					Columns: []*domain.Column{
						{
							Family:    "d",
							Qualifier: "d:row",
							Value:     []byte("madoka"),
							Version:   tm,
						},
					},
				},
				{
					Key: "10",
					Columns: []*domain.Column{
						{
							Family:    "d'",
							Qualifier: "d':row",
							Value:     []byte("madoka"),
							Version:   tm,
						},
					},
				},
			},
		},
		{
			"users",
			bigtable.PrefixRange("4"),
			[]bigtable.ReadOption{
				bigtable.RowFilter(
					bigtable.ChainFilters(
						bigtable.FamilyFilter("^d$"),
						bigtable.LatestNFilter(1),
					),
				),
			},
			[]*domain.Row{
				{
					Key: "4",
					Columns: []*domain.Column{
						{
							Family:    "d",
							Qualifier: "d:row",
							Value:     []byte("anko"),
							Version:   tm.Add(time.Hour),
						},
					},
				},
			},
		},
		{
			"articles",
			bigtable.PrefixRange("3"),
			[]bigtable.ReadOption{},
			[]*domain.Row{
				{
					Key: "3##1",
					Columns: []*domain.Column{
						{
							Family:    "d",
							Qualifier: "d:content",
							Value:     []uint8{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
							Version:   tm,
						},
						{
							Family:    "d",
							Qualifier: "d:title",
							Value:     []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
							Version:   tm,
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		r, err := NewBigtableRepository("test-project", "test-instance")
		assert.NoError(t, err)

		bt, err := r.GetRows(context.Background(), c.table, c.rr, c.opts...)
		assert.NoError(t, err)

		actual := bt.Rows
		assert.Equal(t, c.expect, actual)
	}
}

func TestCount(t *testing.T) {
	loadFixture(t, "testdata/users.yaml")

	cases := []struct {
		table  string
		expect int
	}{
		{"users", 5},
	}
	for _, c := range cases {
		r, err := NewBigtableRepository("test-project", "test-instance")
		assert.NoError(t, err)

		cnt, err := r.Count(context.Background(), c.table)
		assert.NoError(t, err)

		assert.Equal(t, c.expect, cnt)
	}
}

func TestTables(t *testing.T) {
	loadFixture(t, "testdata/users.yaml")
	loadFixture(t, "testdata/articles.yaml")

	cases := []struct {
		expect []string
	}{
		{
			[]string{
				"articles",
				"users",
			},
		},
	}
	for _, c := range cases {
		r, err := NewBigtableRepository("test-project", "test-instance")
		assert.NoError(t, err)

		tbls, err := r.Tables(context.Background())
		assert.NoError(t, err)

		assert.Subset(t, tbls, c.expect)
	}
}

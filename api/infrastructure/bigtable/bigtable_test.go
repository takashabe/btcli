package bigtable

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/takashabe/btcli/api/domain"
)

func TestGet(t *testing.T) {
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
					&domain.Column{
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
					&domain.Column{
						Family:    "d",
						Qualifier: "d:content",
						Value:     []byte("madoka_content"),
						Version:   now,
					},
					&domain.Column{
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

func TestGetRowsWithPrefix(t *testing.T) {
	now := time.Now()
	cases := []struct {
		table  string
		key    string
		expect []*domain.Row
	}{
		{
			"articles",
			"2",
			[]*domain.Row{
				&domain.Row{
					Key: "2##1",
					Columns: []*domain.Column{
						&domain.Column{
							Family:    "d",
							Qualifier: "d:content",
							Value:     []byte("homura_content"),
							Version:   now,
						},
						&domain.Column{
							Family:    "d",
							Qualifier: "d:title",
							Value:     []byte("homura_title"),
							Version:   now,
						},
					},
				},
				&domain.Row{
					Key: "2##2",
					Columns: []*domain.Column{
						&domain.Column{
							Family:    "d",
							Qualifier: "d:content",
							Value:     []byte("homuhomu_content"),
							Version:   now,
						},
						&domain.Column{
							Family:    "d",
							Qualifier: "d:title",
							Value:     []byte("homuhomu_title"),
							Version:   now,
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		r, err := NewBigtableRepository("test-project", "test-instance")
		assert.NoError(t, err)

		bt, err := r.GetRowsWithPrefix(context.Background(), c.table, c.key)
		assert.NoError(t, err)

		// NOTE: hack to version timestamp
		actual := bt.Rows
		for _, row := range actual {
			for _, col := range row.Columns {
				col.Version = now
			}
		}
		assert.Equal(t, c.expect, actual)
	}
}

func TestTables(t *testing.T) {
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

		assert.Equal(t, c.expect, tbls)
	}
}

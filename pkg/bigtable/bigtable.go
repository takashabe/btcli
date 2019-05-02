package bigtable

import (
	"context"
	"io"
	"os"
	"sort"
	"time"

	"cloud.google.com/go/bigtable"
)

//go:generate mockgen --package=bigtable -source=bigtable.go -destination=bigtable_mock.go

// Bigtable entity of the bigtable
type Bigtable struct {
	Table string
	Rows  []*Row
}

// Row represent a row of the table
type Row struct {
	Key     string
	Columns []*Column
}

// Column represent a column of the row
type Column struct {
	Family    string
	Qualifier string
	Value     []byte
	Version   time.Time
}

// Client represent repository of the bigtable
type Client interface {
	OutStream() io.Writer
	ErrStream() io.Writer

	Get(ctx context.Context, table, key string, opts ...bigtable.ReadOption) (*Bigtable, error)
	GetRows(ctx context.Context, table string, rr bigtable.RowRange, opts ...bigtable.ReadOption) (*Bigtable, error)
	Count(ctx context.Context, table string) (int, error)
	Tables(ctx context.Context) ([]string, error)
}

type client struct {
	client      *bigtable.Client
	adminClient *bigtable.AdminClient
	outStream   io.Writer
	errStream   io.Writer
}

// Option functional option pattern for the client.
type Option func(*client)

// NewClient returns initialized client.
func NewClient(project, instance string, opts ...Option) (Client, error) {
	cli, err := getClient(project, instance)
	if err != nil {
		return nil, err
	}
	adminClient, err := getAdminClient(project, instance)
	if err != nil {
		return nil, err
	}
	return &client{
		client:      cli,
		adminClient: adminClient,
	}, nil
}

// WithOutStream settings outStream
func WithOutStream(w io.Writer) Option {
	return func(c *client) {
		c.outStream = w
	}
}

// WithErrStream settings errStream
func WithErrStream(w io.Writer) Option {
	return func(c *client) {
		c.errStream = w
	}
}

func getClient(project, instance string) (*bigtable.Client, error) {
	// TODO: Support options
	return bigtable.NewClient(context.Background(), project, instance)
}

func getAdminClient(project, instance string) (*bigtable.AdminClient, error) {
	// TODO: Support options
	return bigtable.NewAdminClient(context.Background(), project, instance)
}

func (c *client) OutStream() io.Writer {
	if c.outStream == nil {
		return os.Stdout
	}
	return c.outStream
}

func (c *client) ErrStream() io.Writer {
	if c.errStream == nil {
		return os.Stderr
	}
	return c.errStream
}

func (c *client) Get(ctx context.Context, table, key string, opts ...bigtable.ReadOption) (*Bigtable, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tbl := c.client.Open(table)
	row, err := tbl.ReadRow(ctx, key, opts...)
	if err != nil {
		return nil, err
	}
	return &Bigtable{
		Table: table,
		Rows: []*Row{
			readRow(row),
		},
	}, nil
}

func (c *client) GetRows(ctx context.Context, table string, rr bigtable.RowRange, opts ...bigtable.ReadOption) (*Bigtable, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tbl := c.client.Open(table)
	rows := []*Row{}
	err := tbl.ReadRows(ctx, rr, func(row bigtable.Row) bool {
		rows = append(rows, readRow(row))
		return true
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &Bigtable{
		Table: table,
		Rows:  rows,
	}, nil
}

func (c *client) Count(ctx context.Context, table string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tbl := c.client.Open(table)
	cnt := 0
	err := tbl.ReadRows(ctx, bigtable.InfiniteRange(""), func(_ bigtable.Row) bool {
		cnt++
		return true
	}, bigtable.RowFilter(bigtable.StripValueFilter()))
	return cnt, err
}

func readRow(r bigtable.Row) *Row {
	ret := &Row{
		Key:     r.Key(),
		Columns: make([]*Column, 0, len(r)),
	}
	for fam := range r {
		ris := r[fam]
		for _, ri := range ris {
			c := &Column{
				Family:    fam,
				Qualifier: ri.Column,
				Value:     ri.Value,
				Version:   ri.Timestamp.Time(),
			}
			ret.Columns = append(ret.Columns, c)
		}
	}

	sort.Slice(ret.Columns, func(i, j int) bool {
		return ret.Columns[i].Family > ret.Columns[j].Family
	})
	return ret
}

func (c *client) Tables(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tbls, err := c.adminClient.Tables(ctx)
	if err != nil {
		return []string{}, err
	}
	sort.Strings(tbls)
	return tbls, nil
}

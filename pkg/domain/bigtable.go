package domain

import (
	"context"

	"time"

	"cloud.google.com/go/bigtable"
)

//go:generate mockgen --package=domain -source=bigtable.go -destination=bigtable_mock.go

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

// BigtableRepository represent repository of the bigtable
type BigtableRepository interface {
	Get(ctx context.Context, table, key string, opts ...bigtable.ReadOption) (*Bigtable, error)
	GetRows(ctx context.Context, table string, rr bigtable.RowRange, opts ...bigtable.ReadOption) (*Bigtable, error)
	Count(ctx context.Context, table string) (int, error)

	// TODO: Isolation data management client and table management client
	Tables(ctx context.Context) ([]string, error)
}

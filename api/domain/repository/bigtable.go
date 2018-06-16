package repository

//go:generate mockgen --package=repository -source=bigtable.go -destination=bigtable_mock.go

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/takashabe/btcli/api/domain"
)

// Bigtable represent repository of the bigtable
type Bigtable interface {
	Get(ctx context.Context, table, key string) (*domain.Bigtable, error)
	GetRowsWithPrefix(ctx context.Context, table, key string, opts ...bigtable.ReadOption) (*domain.Bigtable, error)
	GetRows(ctx context.Context, table string, rr bigtable.RowRange, opts ...bigtable.ReadOption) (*domain.Bigtable, error)
	Count(ctx context.Context, table string) (int, error)

	// TODO: isolation data management client and table management client
	Tables(ctx context.Context) ([]string, error)
}

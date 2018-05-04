package repository

import (
	"context"

	"github.com/takashabe/btcli/api/domain"
)

// Bigtable represent repository of the bigtable
type Bigtable interface {
	Get(ctx context.Context, table, key string) (*domain.Bigtable, error)
	GetRowsWithPrefix(ctx context.Context, table, key string) (*domain.Bigtable, error)

	// TODO: isolation data management client and table management client
	Tables(ctx context.Context) ([]string, error)
}

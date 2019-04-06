package application

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/takashabe/btcli/pkg/domain"
)

// RowsInteractor provide rows data
type RowsInteractor struct {
	repository domain.BigtableRepository
}

// NewRowsInteractor returns initialized RowsInteractor
func NewRowsInteractor(r domain.BigtableRepository) *RowsInteractor {
	return &RowsInteractor{
		repository: r,
	}
}

// GetRow returns a single row
func (t *RowsInteractor) GetRow(ctx context.Context, table, key string, opts ...bigtable.ReadOption) (*domain.Row, error) {
	tbl, err := t.repository.Get(ctx, table, key, opts...)
	if err != nil {
		return nil, err
	}
	return tbl.Rows[0], nil
}

// GetRows returns rows
func (t *RowsInteractor) GetRows(ctx context.Context, table string, rr bigtable.RowRange, opts ...bigtable.ReadOption) ([]*domain.Row, error) {
	tbl, err := t.repository.GetRows(ctx, table, rr, opts...)
	if err != nil {
		return nil, err
	}
	return tbl.Rows, nil
}

// GetRowCount returns number of the table
func (t *RowsInteractor) GetRowCount(ctx context.Context, table string) (int, error) {
	return t.repository.Count(ctx, table)
}

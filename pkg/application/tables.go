package application

import (
	"context"

	"github.com/takashabe/btcli/pkg/domain"
)

// TableInteractor provide table data
type TableInteractor struct {
	repository domain.BigtableRepository
}

// NewTableInteractor returns initialized TableInteractor
func NewTableInteractor(r domain.BigtableRepository) *TableInteractor {
	return &TableInteractor{
		repository: r,
	}
}

// GetTables returns list table
func (t *TableInteractor) GetTables(ctx context.Context) ([]string, error) {
	return t.repository.Tables(ctx)
}

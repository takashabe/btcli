package application

import (
	"context"

	"github.com/takashabe/btcli/api/domain/repository"
)

// TableInteractor provide table data
type TableInteractor struct {
	repository repository.Bigtable
}

// NewTableInteractor returns initialized TableInteractor
func NewTableInteractor(r repository.Bigtable) *TableInteractor {
	return &TableInteractor{
		repository: r,
	}
}

// GetTables returns list table
func (t *TableInteractor) GetTables(ctx context.Context) ([]string, error) {
	return t.repository.Tables(ctx)
}

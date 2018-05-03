package bigtable

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/takashabe/btcli/api/domain"
	"github.com/takashabe/btcli/api/domain/repository"
)

type bigtableRepository struct {
	client *bigtable.Client
}

// NewBigtableRepository returns initialized bigtableRepository
func NewBigtableRepository(project, instance string) (repository.Bigtable, error) {
	client, err := bigtable.NewClient(context.Background(), project, instance)
	if err != nil {
		return nil, err
	}
	return &bigtableRepository{
		client: client,
	}, nil
}

func (b *bigtableRepository) Get(ctx context.Context, table, key string) (*domain.Bigtable, error) {
	tbl := b.client.Open(table)

	row, err := tbl.ReadRow(ctx, key)
	if err != nil {
		return nil, err
	}
	return &domain.Bigtable{
		Table: table,
		Rows: []*domain.Row{
			readRow(row),
		},
	}, nil
}

func readRow(r bigtable.Row) *domain.Row {
	ret := &domain.Row{
		Key:     r.Key(),
		Columns: make([]*domain.Column, 0, len(r)),
	}
	for fam := range r {
		ris := r[fam]
		for _, ri := range ris {
			c := &domain.Column{
				Family:    fam,
				Qualifier: ri.Column,
				Value:     ri.Value,
				Version:   ri.Timestamp.Time(),
			}
			ret.Columns = append(ret.Columns, c)
		}
	}
	return ret
}

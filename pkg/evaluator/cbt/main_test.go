package cbt

import (
	"testing"

	"cloud.google.com/go/bigtable"
)

func TestMain(m *testing.M) {
	m.Run()
}

func filtersToReadOption(fs ...bigtable.Filter) bigtable.ReadOption {
	return bigtable.RowFilter(bigtable.ChainFilters(fs...))
}

package interfaces

import (
	"os"
	"testing"

	"cloud.google.com/go/bigtable"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func chainFilters(filters ...bigtable.Filter) bigtable.ReadOption {
	return bigtable.RowFilter(bigtable.ChainFilters(filters...))
}

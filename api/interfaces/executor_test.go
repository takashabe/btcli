package interfaces

import (
	"fmt"
	"reflect"
	"testing"

	"cloud.google.com/go/bigtable"
	"github.com/stretchr/testify/assert"
)

func TestRowRange(t *testing.T) {
	cases := []struct {
		input  map[string]string
		expect bigtable.RowRange
	}{
		{
			map[string]string{
				"prefix": "1",
			},
			bigtable.NewRange("1", "2"),
		},
		{
			// TODO: Not supported yet
			map[string]string{
				"start": "1",
				"end":   "2",
			},
			bigtable.NewRange("", ""),
		},
	}
	for _, c := range cases {
		actual, err := rowRange(c.input)
		assert.NoError(t, err)
		assert.Equal(t, c.expect, actual)
	}
}

func TestReadOption(t *testing.T) {
	cases := []struct {
		input   map[string]string
		expects []bigtable.ReadOption
	}{
		{
			map[string]string{
				"count": "1",
			},
			[]bigtable.ReadOption{
				bigtable.LimitRows(0),
			},
		},
		{
			map[string]string{
				"count": "1",
				"regex": "1",
			},
			[]bigtable.ReadOption{
				bigtable.LimitRows(0),
				bigtable.RowFilter(bigtable.RowKeyFilter("")),
			},
		},
	}
	for _, c := range cases {
		actual, err := readOption(c.input)
		assert.NoError(t, err)

		for _, e := range c.expects {
			contain := false
			expectType := reflect.TypeOf(e)
			for _, a := range actual {
				if expectType == reflect.TypeOf(a) {
					contain = true
				}
			}
			if !contain {
				assert.Fail(t, fmt.Sprintf("Expect contan type '%v'", expectType))
			}
		}
	}
}

func TestReadWithOptions(t *testing.T) {
	// TODO: implements used mock
	// ctrl := gomock.NewController(t)
	// mockBtRepo := repository.NewMockBigtable(ctrl)
	// defer ctrl.Finish()
}

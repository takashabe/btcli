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
		input  map[string]string
		expect []bigtable.ReadOption
	}{
		{
			map[string]string{
				"count": "1",
			},
		},
	}
	for _, c := range cases {
		actual, err := readOption(c.input)
		assert.NoError(t, err)

		// TODO: Compare types actual and expect
		fmt.Println(reflect.TypeOf(actual[0]) == reflect.TypeOf(bigtable.LimitRows(0)))
		// _, ok := actual[0].(bigtable.LimitRows)
		// if !ok {
		//   assert.Fail(t, "unmatched readOption type")
		// }
	}
}

func TestReadWithOptions(t *testing.T) {
	// TODO: implements used mock
	// ctrl := gomock.NewController(t)
	// mockBtRepo := repository.NewMockBigtable(ctrl)
	// defer ctrl.Finish()
}

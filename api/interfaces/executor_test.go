package interfaces

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/takashabe/btcli/api/domain/repository"
)

func TestRowRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBtRepo := repository.NewMockBigtable(ctrl)
	defer ctrl.Finish()

	// TODO: Add test case with mockBtRepo

	cases := []struct {
		input  map[string]string
		expect error
	}{}
	for _, c := range cases {
		assert.NoError(t, nil)
	}
}

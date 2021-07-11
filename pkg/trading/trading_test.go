package trading

import (
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRoundToTickSize(t *testing.T) {
	testCases := []struct {
		price    decimal.Decimal
		tickSize decimal.Decimal

		out string
	}{
		{
			price:    decimal.NewFromFloat(1.1234),
			tickSize: decimal.NewFromFloat(0.1),
			out:      "1.1",
		},
		{
			price:    decimal.NewFromFloat(1.123456789),
			tickSize: decimal.NewFromFloat(0.000001),
			out:      "1.123457",
		},
	}
	for i, tt := range testCases {
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tt.out, roundToTickSize(tt.price, tt.tickSize).String())
		})
	}
}

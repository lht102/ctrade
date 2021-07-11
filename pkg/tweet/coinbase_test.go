package tweet

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNewCoinListingPattern(t *testing.T) {
	testCases := []struct {
		in  string
		out bool
	}{
		{
			in:  "Starting today, inbound transfers for FORTH are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin once liquidity conditions are met.",
			out: true,
		},
		{
			in:  "Starting today, inbound transfers for USDT are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 6PM PT on Monday April 26 , if liquidity conditions are met.",
			out: true,
		},
		{
			in:  "Starting today, inbound transfers for SOL are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Monday May 24, if liquidity conditions are met.",
			out: true,
		},
		{
			in:  "Starting today, inbound transfers for DOGE are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Thursday June 3, if liquidity conditions are met.",
			out: true,
		},
		{
			in:  "Starting today, inbound transfers for GTC, MLN & AMP are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Thurs 6/10 if liquidity conditions are met.",
			out: true,
		},
		{
			in:  "Starting today, inbound transfers for DOT are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Wednesday June 16, if liquidity conditions are met.",
			out: true,
		},
		{
			in:  "Inbound transfers for CHZ, KEEP & SHIB are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Thurs 6/17, if liquidity conditions are met.",
			out: true,
		},
		{
			in:  "Inbound transfers for BOND, LPT & QNT are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Wed 6/24, if liquidity conditions are met.",
			out: true,
		},
		{
			in:  "Starting today, inbound transfers for 1INCH, ENJ, NKN & OGN are available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Fri 4/9 if liquidity conditions are met.",
			out: true,
		},
		{
			in:  "Our BOND-USD & LPT-USD order books are now in full-trading mode. Limit, market and stop orders are all now available.",
			out: false,
		},
		{
			in:  "QNT-USD order book will now enter limit-only mode. Limit orders can be placed and cancelled, and matches may occur. Market orders cannot be submitted. The order book will remain in limit-only mode for a minimum of 10 mins.",
			out: false,
		},
		{
			in:  "Our BOND-USD & LPT-USD order books will now enter limit-only mode. Limit orders can be placed and cancelled, and matches may occur. Market orders cannot be submitted. The order book will remain in limit-only mode for a minimum of 10 mins.",
			out: false,
		},
		{
			in:  "Trading on our BOND-USD, LPT-USD & QNT-USD order books is about to begin. Books will now enter post-only mode. Customers can post limit orders but there will be no matches (completed orders). The books will be in post-only mode for a minimum of 1 min.",
			out: false,
		},
		{
			in:  "Our DOT-USD and DOT-BTC order books are now in full-trading mode. Limit, market and stop orders are all now available.",
			out: false,
		},
		{
			in:  "Our DOT-USD & DOT-BTC order books will now enter limit-only mode. Limit orders can be placed and cancelled, and matches may occur. Market orders cannot be submitted.",
			out: false,
		},
		{
			in:  "Our DOT-EUR order book will now enter limit-only mode. Limit orders can be placed and cancelled, and matches may occur. Market orders cannot be submitted. The order book will remain in limit-only mode for a minimum of 10 mins.",
			out: false,
		},
		{
			in:  "Trading on our DOT-USD order book is about to begin. This book will now enter post-only mode. Customers can post limit orders but there will be no matches (completed orders). The books will be in post-only mode for a minimum of 1 min.",
			out: false,
		},
	}
	for i, tt := range testCases {
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tt.out, IsCoinbaseNewCoinListingPattern(tt.in))
		})
	}
}

func TestExtractSymbols(t *testing.T) {
	testCases := []struct {
		in  string
		out []string
	}{
		{
			in:  "Starting today, inbound transfers for FORTH are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin once liquidity conditions are met.",
			out: []string{"FORTH"},
		},
		{
			in:  "Starting today, inbound transfers for SOL are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Monday May 24, if liquidity conditions are met.",
			out: []string{"SOL"},
		},
		{
			in:  "Starting today, inbound transfers for DOGE are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Thursday June 3, if liquidity conditions are met.",
			out: []string{"DOGE"},
		},
		{
			in:  "Starting today, inbound transfers for GTC, MLN & AMP are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Thurs 6/10 if liquidity conditions are met.",
			out: []string{"GTC", "MLN", "AMP"},
		},
		{
			in:  "Starting today, inbound transfers for DOT are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Wednesday June 16, if liquidity conditions are met.",
			out: []string{"DOT"},
		},
		{
			in:  "Inbound transfers for CHZ, KEEP & SHIB are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Thurs 6/17, if liquidity conditions are met.",
			out: []string{"CHZ", "KEEP", "SHIB"},
		},
		{
			in:  "Inbound transfers for BOND, LPT & QNT are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Wed 6/24, if liquidity conditions are met.",
			out: []string{"BOND", "LPT", "QNT"},
		},
		{
			in:  "Starting today, inbound transfers for 1INCH, ENJ, NKN & OGN are available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Fri 4/9 if liquidity conditions are met.",
			out: []string{"1INCH", "ENJ", "NKN", "OGN"},
		},
	}
	for i, tt := range testCases {
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tt.out, extractSymbols(tt.in))
		})
	}
}

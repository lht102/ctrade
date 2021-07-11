package tweet

import (
	"strings"
	"unicode"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/hbollon/go-edlib"
	"github.com/lht102/ctrade/api"
)

const (
	CoinbaseProTwitterUserID = "720487892670410753"
	newCoinListingPattern    = "Starting today, inbound transfers for XXX are now available in the regions where trading is supported. Traders cannot place orders and no orders will be filled. Trading will begin on or after 9AM PT on Mon 1/1 if liquidity conditions are met."
)

func IsCoinbaseNewCoinListingPattern(text string) bool {
	return edlib.JaroWinklerSimilarity(text, newCoinListingPattern) > 0.75 && strings.Contains(text, "transfer")
}

func handleCoinbaseTweetMessage(supportedCoins map[string]struct{}, t *twitter.Tweet, buySignalCh chan api.BuySignal) {
	if t.User.IDStr == CoinbaseProTwitterUserID &&
		!isReply(t) && !isRetweet(t) {
		if IsCoinbaseNewCoinListingPattern(t.Text) {
			symbols := extractSymbols(t.Text)
			for _, s := range symbols {
				if isExist(supportedCoins, s) {
					buySignalCh <- api.BuySignal{
						Symbol: s,
						Source: getTweetURL(t.User.ScreenName, t.IDStr),
					}
				}
			}
		}
	}
}

func extractSymbols(text string) []string {
	words := strings.Fields(text)
	hasSeen := false
	res := []string{}

	for _, w := range words {
		s := strings.Trim(w, "&,.")
		if len(s) == 0 {
			continue
		}

		isValid := isSymbol(s)
		if hasSeen && !isValid {
			break
		}

		if isValid {
			res = append(res, s)
			hasSeen = true
		}
	}

	return res
}

func isSymbol(s string) bool {
	for _, ch := range s {
		if unicode.IsLetter(ch) && !unicode.IsUpper(ch) {
			return false
		}
	}

	return true
}

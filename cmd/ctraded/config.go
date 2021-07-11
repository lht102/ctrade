package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/lht102/ctrade/pkg/trading"
	"github.com/spf13/viper"
	coingecko "github.com/superoo7/go-gecko/v3"
)

const (
	longHTTPTimeout  = 30 * time.Second
	shortHTTPTimeout = 5 * time.Second
)

var (
	errEmptyTwitterAPIKey            = errors.New("empty twitter API key")
	errEmptyTwitterAPISecretKey      = errors.New("empty twitter API secret key")
	errEmptyTwitterAccessToken       = errors.New("empty twitter access token")
	errEmptyTwitterAccessTokenSecret = errors.New("empty twitter access token secret")
	errEmptyBinanceAPIKey            = errors.New("empty binacne api key")
	errEmptyBinanceAPISecretKey      = errors.New("empty binacne api secret key")
)

func getEnv(v *viper.Viper) string {
	return v.GetString("ENV")
}

func getTwitterClientCredentialsConfig(v *viper.Viper) (*oauth1.Config, error) {
	twitterAPIKey := v.GetString("TWITTER_API_KEY")
	if twitterAPIKey == "" {
		return nil, errEmptyTwitterAPIKey
	}

	twitterAPISecretKey := v.GetString("TWITTER_API_SECRET_KEY")
	if twitterAPISecretKey == "" {
		return nil, errEmptyTwitterAPISecretKey
	}

	return oauth1.NewConfig(twitterAPIKey, twitterAPISecretKey), nil
}

func getTwitterAccessTokenConfig(v *viper.Viper) (*oauth1.Token, error) {
	twitterAccessToken := v.GetString("TWITTER_ACCESS_TOKEN")
	if twitterAccessToken == "" {
		return nil, errEmptyTwitterAccessToken
	}

	twitterAccessTokenSecret := v.GetString("TWITTER_ACCESS_TOKEN_SECRET")
	if twitterAccessTokenSecret == "" {
		return nil, errEmptyTwitterAccessTokenSecret
	}

	return oauth1.NewToken(twitterAccessToken, twitterAccessTokenSecret), nil
}

func getBinanceAPIKey(v *viper.Viper) (string, error) {
	binanceAPIKey := v.GetString("BINANCE_API_KEY")
	if binanceAPIKey == "" {
		return "", errEmptyBinanceAPIKey
	}

	return binanceAPIKey, nil
}

func getBinanceAPISecretKey(v *viper.Viper) (string, error) {
	binanceAPISecretKey := v.GetString("BINANCE_API_SECRET_KEY")
	if binanceAPISecretKey == "" {
		return "", errEmptyBinanceAPISecretKey
	}

	return binanceAPISecretKey, nil
}

func getFuturesOptions(v *viper.Viper) []trading.FuturesOption {
	var opts []trading.FuturesOption

	willExecuteOrder := v.GetBool("WILL_EXECUTE_ORDER")
	if willExecuteOrder {
		opts = append(opts, trading.WithWillExecuteOrder(willExecuteOrder))
	}

	leverage := v.GetInt("FUTURES_LEVERAGE")
	if leverage > 0 {
		opts = append(opts, trading.WithLeverage(leverage))
	}

	eachTradeAmountInUSD := v.GetFloat64("FUTURES_EACH_TRADE_AMOUNT_IN_USD")
	if eachTradeAmountInUSD > 0 {
		opts = append(opts, trading.WithEachTradeAmountInUSD(eachTradeAmountInUSD))
	}

	takeProfitPriceChangedPercentage := v.GetFloat64("FUTURES_TAKE_PROFIT_PRICE_CHANGED_PERCENTAGE")
	if takeProfitPriceChangedPercentage != 0 {
		opts = append(opts, trading.WithTakeProfitPriceChangedPercentage(takeProfitPriceChangedPercentage))
	}

	return opts
}

func getSupportedCoins() (map[string]struct{}, error) {
	coingeckoClient := coingecko.NewClient(&http.Client{
		Timeout: longHTTPTimeout,
	})

	coins, err := coingeckoClient.CoinsList()
	if err != nil {
		return nil, fmt.Errorf("coingecko get coins list: %w", err)
	}

	symbolSet := make(map[string]struct{}, len(*coins))
	for _, c := range *coins {
		symbolSet[strings.ToUpper(c.Symbol)] = struct{}{}
	}

	return symbolSet, nil
}

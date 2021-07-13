package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/blendle/zapdriver"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/lht102/ctrade/pkg/trading"
	"github.com/lht102/ctrade/pkg/tweet"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	v := viper.New()
	v.AutomaticEnv()

	if getEnv(v) != "prod" {
		futures.UseTestnet = true
	}

	logger, err := zapdriver.NewProduction()
	if err != nil {
		log.Fatalln("Fail to init logger")
	}

	defer func() {
		_ = logger.Sync()
	}()

	supportedCoins, err := getSupportedCoins()
	if err != nil {
		logger.Fatal("Fail to get supported coins", zap.Error(err))
	}

	binanceAPIKey, err := getBinanceAPIKey(v)
	if err != nil {
		logger.Fatal("Fail to get binance API key", zap.Error(err))
	}

	binanceAPISecretKey, err := getBinanceAPISecretKey(v)
	if err != nil {
		logger.Fatal("Fail to get binance API Secret key", zap.Error(err))
	}

	binanceFuturesClient := binance.NewFuturesClient(binanceAPIKey, binanceAPISecretKey)
	binanceFuturesClient.HTTPClient.Timeout = shortHTTPTimeout

	binanceFuturesManager, err := trading.NewBinanceFuturesManager(
		binanceFuturesClient,
		logger,
		getFuturesOptions(v)...,
	)
	if err != nil {
		logger.Fatal("Fail to init binance futures manager", zap.Error(err))
	}

	ticker := time.NewTicker(updateBinanceExchangeInfoInterval)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if err := binanceFuturesManager.UpdateSupportedSymbols(); err != nil {
				logger.Error("Fail to update supported symbols info", zap.Error(err))
			}
		}
	}()

	twitterAuthCfg, err := getTwitterClientCredentialsConfig(v)
	if err != nil {
		logger.Fatal("Fail to init twitter client credentials config", zap.Error(err))
	}

	twitterAccessTokenCfg, err := getTwitterAccessTokenConfig(v)
	if err != nil {
		logger.Fatal("Fail to init twitter access token config", zap.Error(err))
	}

	httpClient := twitterAuthCfg.Client(context.Background(), twitterAccessTokenCfg)
	httpClient.Timeout = longHTTPTimeout
	twitterClient := twitter.NewClient(httpClient)
	tweetManager := tweet.NewManager(twitterClient, []string{tweet.CoinbaseProTwitterUserID}, supportedCoins)

	buySignalChFromTweet, err := tweetManager.SubscribeBuySignalChannel()
	if err != nil {
		logger.Fatal("Fail to subscribe buy signal", zap.Error(err))
	}

	defer tweetManager.Stop()

	go func() {
		logger.Info("Start listening on buy signal channel from tweets")

		for v := range buySignalChFromTweet {
			logger.Info("Incoming buy signal", zap.String("symbol", v.Symbol), zap.String("source", v.Source))

			if err := binanceFuturesManager.ConsumeBuySignal(v); err != nil {
				logger.Error("Fail to consume buy signal", zap.Error(err))
			}
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	logger.Info("Stop application")
}

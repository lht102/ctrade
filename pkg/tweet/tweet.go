package tweet

import (
	"fmt"
	"sync/atomic"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/lht102/ctrade/api"
)

type Manager struct {
	twitterClient         *twitter.Client
	trackedTwitterUserIDs []string
	supportedCoins        map[string]struct{}

	done            chan struct{}
	usedStreamCount int32
}

func NewManager(twitterClient *twitter.Client, twitterUserIDs []string, supportedCoins map[string]struct{}) *Manager {
	return &Manager{
		twitterClient:         twitterClient,
		trackedTwitterUserIDs: twitterUserIDs,
		supportedCoins:        supportedCoins,
		done:                  make(chan struct{}),
	}
}

func (m *Manager) SubscribeBuySignalChannel() (<-chan api.BuySignal, error) {
	twitterCh, err := m.SubscribeTweetChannel()
	if err != nil {
		return nil, err
	}

	buySignalCh := make(chan api.BuySignal)

	go func() {
		defer close(buySignalCh)

		for t := range twitterCh {
			handleCoinbaseTweetMessage(m.supportedCoins, t, buySignalCh)
		}
	}()

	return buySignalCh, nil
}

func (m *Manager) SubscribeTweetChannel() (<-chan *twitter.Tweet, error) {
	filterParams := &twitter.StreamFilterParams{
		Follow:        m.trackedTwitterUserIDs,
		StallWarnings: twitter.Bool(true),
	}

	stream, err := m.twitterClient.Streams.Filter(filterParams)
	if err != nil {
		return nil, fmt.Errorf("twitter filterd stream: %w", err)
	}

	tweetCh := make(chan *twitter.Tweet)
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		tweetCh <- tweet
	}

	go func() {
		defer close(tweetCh)
		atomic.AddInt32(&m.usedStreamCount, 1)
		demux.HandleChan(stream.Messages)
	}()

	go func() {
		<-m.done
		stream.Stop()
	}()

	return tweetCh, nil
}

func (m *Manager) Stop() {
	for atomic.LoadInt32(&m.usedStreamCount) > 0 {
		atomic.AddInt32(&m.usedStreamCount, -1)
		m.done <- struct{}{}
	}
}

func isExist(set map[string]struct{}, s string) bool {
	_, exist := set[s]

	return exist
}

func isRetweet(t *twitter.Tweet) bool {
	return t.RetweetedStatus != nil
}

func isReply(t *twitter.Tweet) bool {
	return t.InReplyToUserID != 0
}

func getTweetURL(screenName string, tweetID string) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", screenName, tweetID)
}

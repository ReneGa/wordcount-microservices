package service

import (
	"sync"
	"time"

	"github.com/ReneGa/tweetcount-microservices/windower/domain"
	"github.com/ReneGa/tweetcount-microservices/windower/gateway"
)

type Window struct {
	sync.Mutex
	tweetWordCountsGateway gateway.TweetWordCounts
	searchesGateway        gateway.Searches
	forSearch              map[string]*domain.Window
}

func NewWindow(
	tweetWordCountsGateway gateway.TweetWordCounts,
	searchesGateway gateway.Searches,
) *Window {
	return &Window{
		tweetWordCountsGateway: tweetWordCountsGateway,
		searchesGateway:        searchesGateway,
		forSearch:              map[string]*domain.Window{},
	}
}

func (w *Window) Totals(searchID string) domain.WordCount {
	w.Lock()
	window, ok := w.forSearch[searchID]
	w.Unlock()
	if ok {
		w.Lock()
		defer w.Unlock()
		totals := domain.WordCount{}
		for word, count := range window.Totals {
			totals[word] = count
		}
		return totals
	}

	// Fetch search from searches service
	search := w.searchesGateway.ForID(searchID)

	// Create new window
	window = domain.NewWindow(search.WindowLengthSeconds, 16384) // TODO: Read from config

	// Write into forSearch map
	w.Lock()
	w.forSearch[searchID] = window
	w.Unlock()

	go func() {
		tweetWordCounts := w.tweetWordCountsGateway.TweetWordCounts(search.Query)
		for tweetWordCount := range tweetWordCounts.Data {
			w.Lock()
			window.Enqueue(tweetWordCount)
			window.Trim(time.Now())
			w.Unlock()
		}
	}()

	return domain.WordCount{}
}

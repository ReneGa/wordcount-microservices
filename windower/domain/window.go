package domain

import "time"

// Window is a bounded sliding time window of TweetWordCounts.
type Window struct {
	tweetWordCounts []TweetWordCount
	Totals          WordCount
	MaxCount        int
	LengthSeconds   int
}

// NewWindow allocates a new Window with the given size and bound.
func NewWindow(lengthSeconds, maxCount int) *Window {
	return &Window{
		tweetWordCounts: make([]TweetWordCount, 0, 0),
		Totals:          WordCount{},
		MaxCount:        maxCount,
		LengthSeconds:   lengthSeconds,
	}
}

// Enqueue adds a TweetWordCount to this Window
func (w *Window) Enqueue(tweetWordCount TweetWordCount) {
	for word, count := range tweetWordCount.WordCount {
		w.Totals[word] += count
	}
	w.tweetWordCounts = append(w.tweetWordCounts, tweetWordCount)
}

// Dequeue removes the oldest TweetWordCount from this Window
func (w *Window) Dequeue() {
	if len(w.tweetWordCounts) == 0 {
		return
	}
	tweetWordCount := w.tweetWordCounts[0]
	w.tweetWordCounts = w.tweetWordCounts[1:]
	for word, count := range tweetWordCount.WordCount {
		total := w.Totals[word]
		newTotal := total - count
		if newTotal == 0 {
			delete(w.Totals, word)
		} else {
			w.Totals[word] = newTotal
		}
	}
}

// Trim removes tweets that exceed the bound or are too old
func (w *Window) Trim(now time.Time) {
	for len(w.tweetWordCounts) > w.MaxCount {
		w.Dequeue()
	}
	for w.lastTweetIsTooOld(now) {
		w.Dequeue()
	}
}

func (w *Window) lastTweetIsTooOld(now time.Time) bool {
	if len(w.tweetWordCounts) == 0 {
		return false
	}
	tweetWordCount := w.tweetWordCounts[0]
	windowLength := time.Second * time.Duration(w.LengthSeconds)
	windowBoundary := now.Add(-windowLength)
	return tweetWordCount.TweetTime.Before(windowBoundary)
}

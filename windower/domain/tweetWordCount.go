package domain

import "time"

type WordCount map[string]int

type TweetWordCount struct {
	WordCount     WordCount
	TweetID       string
	TweetTime     time.Time
	TweetLanguage string
}

// TweetWordCounts are a stoppable stream of TweetWordCount-s
type TweetWordCounts struct {
	Data <-chan TweetWordCount
	Stop chan<- bool
}

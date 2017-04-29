package domain

import "time"

type WordCount map[string]int

type TweetWordCount struct {
	WordCount     WordCount
	TweetID       string
	TweetTime     time.Time
	TweetLanguage string
}

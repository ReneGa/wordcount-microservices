package domain

import "time"

// TweetWords represents the words contained in a tweet
type TweetWords struct {
	Words     []string
	TweetID   string
	TweetTime time.Time
}

// TweetsWords is a stoppable stream of TweetWords
type TweetsWords struct {
	Data <-chan TweetWords
	Stop chan<- bool
}

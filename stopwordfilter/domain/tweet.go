package domain

import "time"

// Tweet is the domain object representing a tweet
type Tweet struct {
	ID   string
	Text string
	Time time.Time
}

// TweetWords represents the words contained in a tweet
type TweetWords struct {
	Words     []string
	TweetID   string
	TweetTime time.Time
}

// Tweets are a stoppable stream of tweets
type Tweets struct {
	Data <-chan Tweet
	Stop chan<- bool
}

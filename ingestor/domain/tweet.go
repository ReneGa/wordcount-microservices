package domain

import "time"

// Tweet is the textual content of a tweet
type Tweet struct {
	Text string
	Time time.Time
}

// Tweets are a stoppable stream of tweets
type Tweets struct {
	Tweets <-chan Tweet
	Stop   chan<- struct{}
}

// StopTweetsMessage stops a Tweets stream
type StopTweetsMessage struct{}

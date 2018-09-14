package domain

import "time"

// Tweet is the textual content of a tweet
type Tweet struct {
	ID       string
	Text     string
	Language string
	Time     time.Time
}

// Tweets are a stoppable stream of tweets
type Tweets struct {
	Data <-chan Tweet
	Stop chan<- bool
}

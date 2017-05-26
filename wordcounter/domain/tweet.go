package domain

import "time"
import "strings"

// Tweet is the domain object representing a tweet
type Tweet struct {
	ID       string
	Text     string
	Language string
	Time     time.Time
}

const cutset = ",.:;-_!?()[]{}…\"\n\t'"

func (t Tweet) WordCount() TweetWordCount {
	words := strings.Split(t.Text, " ")
	wordCount := WordCount{}
	for _, word := range words {
		word = strings.Trim(word, cutset)
		if word != "" {
			wordCount[word]++
		}
	}
	return TweetWordCount{
		WordCount:     wordCount,
		TweetID:       t.ID,
		TweetTime:     t.Time,
		TweetLanguage: t.Language,
	}
}

// Tweets are a stoppable stream of tweets
type Tweets struct {
	Data <-chan Tweet
	Stop chan<- bool
}

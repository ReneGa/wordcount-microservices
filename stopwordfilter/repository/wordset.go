package repository

import (
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/datamapper"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/domain"
)

// WordSet is a word set repository
type WordSet interface {
	Get() domain.WordSet
}

// NewWordSet creates a new WordSet repository
func NewWordSet(dataMapper datamapper.WordSet) WordSet {
	return dataMapper
}

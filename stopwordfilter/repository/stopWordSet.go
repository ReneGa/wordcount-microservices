package repository

import (
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/domain"
)

// StopWordSet is a word set repository
type StopWordSet interface {
	Get(ID string) domain.WordSet
	List() []string
}

package repository

import (
	"sync"

	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/datamapper"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/domain"
)

// StopWordSet is a word set repository
type StopWordSet interface {
	Get(ID string) domain.WordSet
	List() []string
}

// NewStopWordSet creates a new StopWordSet repository
func NewStopWordSet(datamapper datamapper.StopWordSet) StopWordSet {
	return &stopWordSet{
		cache:      map[string]domain.WordSet{},
		datamapper: datamapper,
	}
}

type stopWordSet struct {
	sync.RWMutex
	cache      map[string]domain.WordSet
	datamapper datamapper.StopWordSet
}

func (s *stopWordSet) List() []string {
	return s.datamapper.List()
}

func (s *stopWordSet) Get(ID string) domain.WordSet {
	s.RLock()
	stopWordSet, ok := s.cache[ID]
	s.RUnlock()
	if ok {
		return stopWordSet
	}

	stopWordSet = s.datamapper.Get(ID)

	s.Lock()
	s.cache[ID] = stopWordSet
	s.Unlock()

	return stopWordSet
}

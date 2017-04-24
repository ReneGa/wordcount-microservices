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
	return &wordSet{
		cache:      map[string]domain.WordSet{},
		datamapper: datamapper,
	}
}

type wordSet struct {
	sync.RWMutex
	cache      map[string]domain.WordSet
	datamapper datamapper.StopWordSet
}

func (w *wordSet) List() []string {
	return w.datamapper.List()
}

func (w *wordSet) Get(ID string) domain.WordSet {
	w.RLock()
	wordSet, ok := w.cache[ID]
	w.RUnlock()
	if ok {
		return wordSet
	}

	wordSet = w.datamapper.Get(ID)

	w.Lock()
	w.cache[ID] = wordSet
	w.Unlock()

	return wordSet
}

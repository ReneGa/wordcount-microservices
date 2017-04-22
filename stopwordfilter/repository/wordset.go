package repository

import (
	"sync"

	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/datamapper"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/domain"
)

// WordSet is a word set repository
type WordSet interface {
	Get(ID string) domain.WordSet
	List() []string
}

// NewWordSet creates a new WordSet repository
func NewWordSet(datamapper datamapper.WordSet) WordSet {
	return &wordSet{
		cache:      map[string]domain.WordSet{},
		datamapper: datamapper,
	}
}

type wordSet struct {
	sync.RWMutex
	cache      map[string]domain.WordSet
	datamapper datamapper.WordSet
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

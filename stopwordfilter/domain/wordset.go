package domain

import (
	"strings"
)

// WordSet is a set of words
type WordSet map[string]struct{}

// Add adds a word to the word set
func (w WordSet) Add(word string) {
	w[word] = struct{}{}
}

// Contains returns true if the given word is contained in the word set,
// and false otherwise.
func (w WordSet) Contains(word string) bool {
	_, ok := w[strings.ToLower(word)]
	return ok
}

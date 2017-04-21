package domain

import (
	"bufio"
	"os"
	"strings"
)

// WordSet is a set of words
type WordSet map[string]struct{}

// LoadWordSet loads a word set from a file
// FIXME: move to datamapper
func LoadWordSet(fileName string) WordSet {
	wordSet := WordSet{}
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordSet.Add(scanner.Text())
	}
	return wordSet
}

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

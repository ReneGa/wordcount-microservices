package datamapper

import (
	"bufio"
	"os"

	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/domain"
)

// WordSet is a word set datamapper
type WordSet interface {
	Get() domain.WordSet
}

type wordSet struct {
	fileName string
}

// NewWordSet creates a new WordSet flat-file datamapper
func NewWordSet(fileName string) WordSet {
	return &wordSet{fileName}
}

// Load loads a word set from a file
func (w *wordSet) Get() domain.WordSet {
	wordSet := domain.WordSet{}
	file, err := os.Open(w.fileName)
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

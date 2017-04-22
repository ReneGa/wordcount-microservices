package datamapper

import (
	"bufio"
	"os"
	"path/filepath"

	"strings"

	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/domain"
)

const wordSetFileSuffix = ".txt"

// WordSet is a word set datamapper
type WordSet interface {
	Get(ID string) domain.WordSet
	List() []string
}

type wordSet struct {
	wordSetsDirectory string
}

// NewWordSet creates a new WordSet flat-file datamapper
func NewWordSet(wordSetsDirectory string) WordSet {
	return &wordSet{wordSetsDirectory}
}

// Load loads a word set from a file
func (w *wordSet) Get(ID string) domain.WordSet {
	wordSet := domain.WordSet{}
	file, err := os.Open(filepath.Join(w.wordSetsDirectory, ID, wordSetFileSuffix))
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

func (w *wordSet) List() []string {
	IDs := make([]string, 0, 0)
	filepath.Walk(w.wordSetsDirectory, func(path string, info os.FileInfo, err error) error {
		ID, _ := filepath.Rel(w.wordSetsDirectory, path)
		ID = strings.TrimSuffix(ID, wordSetFileSuffix)
		IDs = append(IDs, ID)
		return nil
	})
	return IDs
}

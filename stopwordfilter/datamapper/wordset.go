package datamapper

import (
	"bufio"
	"io/ioutil"
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
	file, err := os.Open(filepath.Join(w.wordSetsDirectory, ID+wordSetFileSuffix))
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
	files, err := ioutil.ReadDir(w.wordSetsDirectory)
	if err != nil {
		panic(err)
	}
	IDs := make([]string, len(files))
	for i, file := range files {
		ID := strings.TrimSuffix(file.Name(), wordSetFileSuffix)
		IDs[i] = ID
	}
	return IDs
}

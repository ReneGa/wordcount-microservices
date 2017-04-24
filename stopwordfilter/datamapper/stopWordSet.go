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

// StopWordSet is a stopword set datamapper
type StopWordSet interface {
	Get(ID string) domain.WordSet
	List() []string
}

type stopWordSet struct {
	wordSetsDirectory string
}

// NewStopWordSet creates a new WordSet flat-file datamapper
func NewStopWordSet(wordSetsDirectory string) StopWordSet {
	return &stopWordSet{wordSetsDirectory}
}

// Load loads a word set from a file
func (w *stopWordSet) Get(ID string) domain.WordSet {
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

func (w *stopWordSet) List() []string {
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

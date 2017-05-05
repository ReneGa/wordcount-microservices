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

// NewDirectoryStopWordSet creates DirectoryStopWordSet
func NewDirectoryStopWordSet(wordSetsDirectory string) *DirectoryStopWordSet {
	d := DirectoryStopWordSet{WordSetsDirectory: wordSetsDirectory}
	d.wordSets = make(map[string]domain.WordSet)
	for _, wordSetID := range d.listFilesOnDisk() {
		d.wordSets[wordSetID] = d.load(wordSetID)
	}
	return &d
}

// DirectoryStopWordSet reads stop words from the given directory
type DirectoryStopWordSet struct {
	WordSetsDirectory string
	wordSets          map[string]domain.WordSet
}

// Get loads a word set from a file
func (d *DirectoryStopWordSet) Get(ID string) domain.WordSet {
	return d.wordSets[ID]
}

func (d *DirectoryStopWordSet) load(ID string) domain.WordSet {
	wordSet := domain.WordSet{}
	file, err := os.Open(filepath.Join(d.WordSetsDirectory, ID+wordSetFileSuffix))
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

// List all available StopwordSets
func (d *DirectoryStopWordSet) List() []string {
	wordSetIDs := make([]string, len(d.wordSets))
	i := 0
	for wordSetID := range d.wordSets {
		wordSetIDs[i] = wordSetID
		i++
	}
	return wordSetIDs
}

func (d *DirectoryStopWordSet) listFilesOnDisk() []string {
	files, err := ioutil.ReadDir(d.WordSetsDirectory)
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

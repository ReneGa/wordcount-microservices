package domain

import "strings"

// removeWords removes the given stop words from the given word slice
func removeWords(wordsToRemove WordSet, words []string) []string {
	filteredWords := make([]string, 0, len(words))
	for _, word := range words {
		if !wordsToRemove.Contains(word) {
			filteredWords = append(filteredWords, word)
		}
	}
	return filteredWords
}

// FilterStopWords removes stop words from a given tweet and returns
// the remaining words in their original order.
func FilterStopWords(stopWords WordSet, tweet Tweet) Tweet {
	words := strings.Split(tweet.Text, " ")
	return Tweet{
		Text:     strings.Join(removeWords(stopWords, words), " "),
		ID:       tweet.ID,
		Time:     tweet.Time,
		Language: tweet.Language,
	}
}

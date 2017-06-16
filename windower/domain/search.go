package domain

// Search is a windowed search query
type Search struct {
	Query               string
	WindowLengthSeconds int
}

// SearchID is the id of a persistable search
type SearchID string

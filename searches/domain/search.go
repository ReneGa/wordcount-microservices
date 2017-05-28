package domain

// Search is a windowed search query
type Search struct {
	ID                  int64
	Query               string
	WindowLengthSeconds int
}

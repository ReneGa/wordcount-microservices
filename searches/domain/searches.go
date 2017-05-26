package domain

// Search is a windowed search query
type Search struct {
	Query               string
	WindowLengthSeconds int
}

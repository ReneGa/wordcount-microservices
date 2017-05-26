package repository

import "github.com/ReneGa/tweetcount-microservices/searches/domain"

// Searches is the repository for Searches
type Searches struct{}

// GetAll returns all stored Searches
func (s *Searches) GetAll() ([]domain.Search, error) {
	return nil, nil // Todo
}

// Save stores Searches to the database
func (s *Searches) Save(*domain.Search) (*domain.Search, error) {
	return nil, nil // Todo
}

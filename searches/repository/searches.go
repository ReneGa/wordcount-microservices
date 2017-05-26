package repository

import (
	"database/sql"
	"fmt"

	"github.com/ReneGa/tweetcount-microservices/searches/domain"

	// blank import
	_ "github.com/mattn/go-sqlite3"
)

// Searches is the repository for Searches
type Searches struct {
	DB *sql.DB
}

// GetAll returns all stored Searches
func (s *Searches) GetAll() ([]domain.Search, error) {
	return nil, nil // Todo
}

// Save stores Searches to the database
func (s *Searches) Save(search *domain.Search) (*domain.Search, error) {

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(
		fmt.Sprintf("insert into searches(query, windowSeconds) values(%s, %d)", search.Query, search.WindowLengthSeconds))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		return nil, err
	}
	tx.Commit()

	search.ID, err = result.LastInsertId()

	return search, err
}

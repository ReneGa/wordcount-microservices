package datamapper

import (
	"database/sql"
	"errors"

	"github.com/ReneGa/tweetcount-microservices/searches/domain"

	// blank import
	_ "github.com/mattn/go-sqlite3"
)

// ErrSearchNotFound is returned by Get when no search
// with the given ID is found
var ErrSearchNotFound = errors.New("search not found")

// Searches is the repository for Searches
type Searches struct {
	DB *sql.DB
}

func (s *Searches) Init() error {
	_, err := s.DB.Exec(`
	create table if not exists searches (
		id integer not null primary key,
		query text not null,
		windowSeconds integer not null
	);`)
	return err
}

func (s *Searches) Get(ID string) (*domain.Search, error) {
	var IDfromDB int64
	var Query string
	var WindowLengthSeconds int
	err := s.DB.QueryRow("select id, query, windowSeconds from searches where id = ?", ID).Scan(&IDfromDB, &Query, &WindowLengthSeconds)
	if err == sql.ErrNoRows {
		return nil, ErrSearchNotFound
	}
	if err != nil {
		return nil, err
	}
	return &domain.Search{
		ID:                  IDfromDB,
		Query:               Query,
		WindowLengthSeconds: WindowLengthSeconds,
	}, nil
}

// GetAll returns all stored Searches
func (s *Searches) GetAll() ([]domain.Search, error) {
	stmt, err := s.DB.Prepare("select id, query, windowSeconds from searches")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	defer rows.Close()
	searches := make([]domain.Search, 0, 0)
	for rows.Next() {
		var ID int64
		var Query string
		var WindowLengthSeconds int
		err := rows.Scan(&ID, &Query, &WindowLengthSeconds)
		if err != nil {
			return nil, err
		}
		search := domain.Search{
			ID:                  ID,
			Query:               Query,
			WindowLengthSeconds: WindowLengthSeconds,
		}
		searches = append(searches, search)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return searches, nil
}

// Save stores Searches to the database
func (s *Searches) Save(search *domain.Search) (*domain.Search, error) {

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare("insert into searches(query, windowSeconds) values (?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(search.Query, search.WindowLengthSeconds)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	search.ID, err = result.LastInsertId()

	return search, err
}

package gateway

import (
	"errors"
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/windower/domain"
)

// Searches is a gateway to a Search service
type Searches interface {
	ForID(ID domain.SearchID) (*domain.Search, error)
}

// HTTPSearches is a gateway to a HTTP Search service
type HTTPSearches struct {
	Client *http.Client
	URL    string
}

// ErrSearchNotFound is the error indicating that the given search was not found
var ErrSearchNotFound = errors.New("search not found")

// ForID returns the search for the given ID
func (s *HTTPSearches) ForID(ID string) (*domain.Search, error) {
	url := fmt.Sprintf("%s/%s", s.URL, ID)
	res, err := s.Client.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusNotFound {
		return nil, ErrSearchNotFound
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("searches: unexpected response status code %d", res.StatusCode)
	}
	jd := json.NewDecoder(res.Body)
	var search domain.Search
	err = jd.Decode(&search)
	if err != nil {
		return nil, err
	}
	return &search, nil
}

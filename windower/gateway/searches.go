package gateway

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/windower/domain"
)

type Searches interface {
	ForID(ID string) domain.Search
}

type HTTPSearches struct {
	Client *http.Client
	URL    string
}

func (s *HTTPSearches) ForID(ID string) domain.Search {
	url := fmt.Sprintf("%s/%s", s.URL, ID)
	res, err := s.Client.Get(url)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("searches: unexpected response status code %d", res.StatusCode))
	}
	jd := json.NewDecoder(res.Body)
	var search domain.Search
	err = jd.Decode(&search)
	if err != nil {
		panic(err)
	}
	return search
}

package gateway

import "github.com/ReneGa/tweetcount-microservices/windower/domain"

type Searches interface {
	ForID(ID string) domain.Search
}

type FixedSearches struct {
	domain.Search
}

func (f *FixedSearches) ForID(ID string) domain.Search {
	return domain.Search{
		Query:               f.Query,
		WindowLengthSeconds: f.WindowLengthSeconds,
	}
}

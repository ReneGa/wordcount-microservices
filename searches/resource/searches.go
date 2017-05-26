package resource

import (
	"net/http"

	"encoding/json"

	"strconv"

	"github.com/ReneGa/tweetcount-microservices/searches/domain"
	"github.com/julienschmidt/httprouter"
)

// Search is a resource serving window Search
type Search struct {
	SearchesRepository *repository.Searches
}

// GetAll retrieves all persisted searches
func (s *Search) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	searches := s.SearchesRepository.GetAll()
	je := json.NewEncoder(w)
	je.Encode(searches)
}

// Post creates a new search
func (s *Search) Post(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := r.URL.Query().Get("q")
	seconds, err := strconv.Atoi(r.URL.Query().Get("seconds"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	newSearch := domain.Search{
		Query:               query,
		WindowLengthSeconds: seconds,
	}
	search := s.SearchesRepository.Save(newSearch)
	je := json.NewEncoder(w)
	je.Encode(search)
}

package resource

import (
	"net/http"

	"encoding/json"

	"strconv"

	"github.com/ReneGa/tweetcount-microservices/searches/domain"
	"github.com/ReneGa/tweetcount-microservices/searches/repository"
	"github.com/julienschmidt/httprouter"
)

// Searches is a resource serving window Search
type Searches struct {
	SearchesRepository *repository.Searches
}

func (s *Searches) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	searches, err := s.SearchesRepository.Get(p.ByName("searchID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	je := json.NewEncoder(w)
	je.Encode(searches)
}

// GetAll retrieves all persisted searches
func (s *Searches) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	searches, err := s.SearchesRepository.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	je := json.NewEncoder(w)
	je.Encode(searches)
}

// Create creates a new search
func (s *Searches) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
	search, err := s.SearchesRepository.Save(&newSearch)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	je := json.NewEncoder(w)
	je.Encode(search)
}

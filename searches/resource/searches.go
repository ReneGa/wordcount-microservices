package resource

import (
	"net/http"

	"encoding/json"

	"strconv"

	"github.com/ReneGa/tweetcount-microservices/searches/datamapper"
	"github.com/ReneGa/tweetcount-microservices/searches/domain"
	"github.com/julienschmidt/httprouter"
)

// Searches is a resource serving window Search
type Searches struct {
	SearchesDataMapper *datamapper.Searches
}

func (s *Searches) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	searches, err := s.SearchesDataMapper.Get(p.ByName("searchID"))
	if err == datamapper.ErrSearchNotFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}
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
	searches, err := s.SearchesDataMapper.GetAll()
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
	search, err := s.SearchesDataMapper.Save(&newSearch)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	je := json.NewEncoder(w)
	je.Encode(search)
}

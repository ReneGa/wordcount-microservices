package resource

import (
	"net/http"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/windower/domain"
	"github.com/ReneGa/tweetcount-microservices/windower/service"
	"github.com/julienschmidt/httprouter"
)

// Totals is a resource serving window totals
type Totals struct {
	Service *service.Window
}

// GET writes the current window state to the response
func (t *Totals) GET(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	searchID := domain.SearchID(p.ByName("searchID"))
	totals, err := t.Service.Totals(searchID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	je := json.NewEncoder(w)
	je.Encode(totals)
}

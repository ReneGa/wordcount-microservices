package resource

import (
	"net/http"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/windower/service"
	"github.com/julienschmidt/httprouter"
)

// Totals is a resource serving window totals
type Totals struct {
	Service *service.Window
}

// GET writes the current window state to the response
func (t *Totals) GET(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	searchID := p.ByName("searchID")
	totals := t.Service.Totals(searchID)
	je := json.NewEncoder(w)
	je.Encode(totals)
}

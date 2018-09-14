package resource

import (
	"net/http"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/generator/domain"
	"github.com/julienschmidt/httprouter"
)

// Tweets is the tweets resource
type Tweets struct {
	Generator *domain.TweetsGenerator
}

// GET generates tweets for a given search query
func (t *Tweets) GET(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := r.URL.Query().Get("q")

	// open tweets stream
	tweets := t.Generator.Tweets(query)

	// connection close should stop tweet stream
	go func() {
		tweets.Stop <- <-w.(http.CloseNotifier).CloseNotify()
	}()

	// write values to response
	for tweet := range tweets.Data {
		je := json.NewEncoder(w)
		je.Encode(tweet)
		w.(http.Flusher).Flush()
	}
}

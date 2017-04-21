package resource

import (
	"net/http"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/ingestor/gateway"
	"github.com/julienschmidt/httprouter"
)

// Tweets is the tweets resource
type Tweets struct {
	Gateway gateway.Twitter
}

// GET streams tweets from Twitter for a given search query
func (t *Tweets) GET(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := r.URL.Query().Get("q")

	// open tweets stream
	tweets := t.Gateway.Tweets(query)

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

package resource

import (
	"net/http"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/wordcounter/gateway"
	"github.com/julienschmidt/httprouter"
)

// WordCounts is a resource serving word counts of tweets
type WordCounts struct {
	Gateway gateway.Tweets
}

// GET writes a stream of tweet word counts to the response
func (t *WordCounts) GET(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := r.URL.Query().Get("q")
	offset := r.URL.Query().Get("t")

	tweets := t.Gateway.Tweets(query, offset)

	go func() {
		tweets.Stop <- <-w.(http.CloseNotifier).CloseNotify()
	}()

	je := json.NewEncoder(w)
	for tweet := range tweets.Data {
		je.Encode(tweet.WordCount())
		w.(http.Flusher).Flush()
	}
}

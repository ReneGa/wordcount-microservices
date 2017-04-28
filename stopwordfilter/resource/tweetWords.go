package resource

import (
	"net/http"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/gateway"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/service"
	"github.com/julienschmidt/httprouter"
)

// Tweets is a resource serving tweets passed through a stopword filter
type Tweets struct {
	Gateway gateway.Tweets
	Service service.StopWordFilter
}

// GET writes a stream of filtered tweets to the response
func (t *Tweets) GET(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := r.URL.Query().Get("q")
	tweets := t.Gateway.Tweets(query)
	filteredTweets := t.Service.Filter(tweets)

	go func() {
		filteredTweets.Stop <- <-w.(http.CloseNotifier).CloseNotify()
	}()

	je := json.NewEncoder(w)
	for tweetWords := range filteredTweets.Data {
		je.Encode(tweetWords)
		w.(http.Flusher).Flush()
	}
}

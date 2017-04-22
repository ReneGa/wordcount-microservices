package resource

import (
	"net/http"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/gateway"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/service"
	"github.com/julienschmidt/httprouter"
)

// TweetsWords is a resource serving tweets passed through a stopword filter
type TweetsWords struct {
	Gateway gateway.Tweets
	Service service.StopWordFilter
}

// GET writes a stream of filtered tweets to the response
func (t *TweetsWords) GET(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := r.URL.Query().Get("q")
	tweets := t.Gateway.Tweets(query)
	filteredTweetsWords := t.Service.TweetsWords(tweets)

	go func() {
		filteredTweetsWords.Stop <- <-w.(http.CloseNotifier).CloseNotify()
	}()

	je := json.NewEncoder(w)
	for tweetWords := range filteredTweetsWords.Data {
		je.Encode(tweetWords)
		w.(http.Flusher).Flush()
	}
}

package resource

import (
	"net/http"
	"time"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/recorder/service"
	"github.com/julienschmidt/httprouter"
)

// Tweets is the tweets resource
type Tweets struct {
	Service service.Tweets
}

// GET streams tweets from Twitter for a given search query
func (t *Tweets) GET(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := r.URL.Query().Get("q")
	offsetString := r.URL.Query().Get("t")

	// parse start time offset
	offset, err := time.ParseDuration(offsetString)

	// if given, subtract start time offset from `now`
	startTime := time.Now()
	if err == nil {
		startTime = startTime.Add(-offset)
	}

	// open tweets stream
	tweets := t.Service.Tweets(query, startTime)

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

package domain

import (
	"log"
	"math/big"
	"sync"
	"time"
)

type TweetsGenerator struct {
	sync.Mutex
	ID       *big.Int
	Delay    time.Duration
	Language string
}

var one = big.NewInt(1)
var zero = big.NewInt(0)

func (t *TweetsGenerator) initID() {
	t.Lock()
	if t.ID == nil {
		t.ID = zero
	}
	t.Unlock()
}

func (t *TweetsGenerator) incrementID() {
	t.Lock()
	t.ID.Add(t.ID, one)
	t.Unlock()
}

func (t *TweetsGenerator) generateTweet(query string) Tweet {
	return Tweet{
		Text:     query,
		ID:       t.ID.String(),
		Language: t.Language,
		Time:     time.Now(),
	}
}

func (t *TweetsGenerator) Tweets(query string) Tweets {
	data := make(chan Tweet)
	stop := make(chan bool)
	t.initID()

	go func() {
		for {
			tweet := t.generateTweet(query)
			select {
			case data <- tweet:
				log.Println("Wrote tweet", tweet)
				t.incrementID()
				time.Sleep(t.Delay)
			case <-stop:
				return
			}
		}
	}()

	return Tweets{
		Data: data,
		Stop: stop,
	}
}

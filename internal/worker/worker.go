package worker

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Pool struct {
	client    Doer
	numWorker int
	rate      int
}

func NewPool(client Doer, numWorker, rate int) *Pool {
	return &Pool{
		client:    client,
		numWorker: numWorker,
		rate:      rate,
	}
}

func (p *Pool) Run(reqs <-chan *http.Request) {
	var wg sync.WaitGroup

	var ticker *time.Ticker
	if p.rate > 0 {
		ticker = time.NewTicker(time.Second / time.Duration(p.rate))
	}

	for i := 0; i < p.numWorker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for req := range reqs {
				if ticker != nil {
					<-ticker.C
				}
				resp, err := p.client.Do(req)
				if err != nil {
					fmt.Printf("request error: %v\n", err)
					continue
				}
				if resp != nil && resp.Body != nil {
					resp.Body.Close()
				}
				fmt.Printf("status code: %d\n", resp.StatusCode)
			}
		}()
	}
	wg.Wait()
	if ticker != nil {
		ticker.Stop()
	}
}
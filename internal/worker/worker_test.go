package worker

import (
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type mockDoer struct {
	doFunc func(*http.Request) (*http.Response, error)
}

func (m *mockDoer) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func TestWorker(t *testing.T) {
	var count int32
	mockClient := &mockDoer{
		doFunc: func(req *http.Request) (*http.Response, error) {
			atomic.AddInt32(&count, 1)
			return &http.Response{StatusCode: http.StatusOK}, nil
		},
	}

	pool := NewPool(mockClient, 2, 0)

	reqs := make(chan *http.Request, 5)
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)
		reqs <- req
	}
	close(reqs)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		pool.Run(reqs)
	}()
	wg.Wait()

	if count != 5 {
		t.Errorf("expected 5 requests to be processed, but got %d", count)
	}
}

func TestWorker_WithRateLimit(t *testing.T) {
	var count int32
	mockClient := &mockDoer{
		doFunc: func(req *http.Request) (*http.Response, error) {
			atomic.AddInt32(&count, 1)
			return &http.Response{StatusCode: http.StatusOK}, nil
		},
	}

	// 1秒あたり2リクエストに制限
	pool := NewPool(mockClient, 1, 2)

	reqs := make(chan *http.Request, 5)
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)
		reqs <- req
	}
	close(reqs)

	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		pool.Run(reqs)
	}()
	wg.Wait()
	duration := time.Since(start)

	if count != 5 {
		t.Errorf("expected 5 requests to be processed, but got %d", count)
	}
	// 5リクエストを2req/secで処理すると、最低でも2秒はかかるはず
	if duration < 2*time.Second {
		t.Errorf("expected duration to be at least 2s, but got %v", duration)
	}
}
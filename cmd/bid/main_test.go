package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestClient is fake client which returns predefined results
// TestClient implements Clienter
type TestClient struct {
	Sources map[string]*TestResponse
}

// TestResponse represents test examples of source response
type TestResponse struct {
	Timeout int           // response time from source in ms
	Code    int           // status code of response
	Body    io.ReadCloser // reader for body
}

// TestBody implements io.ReadCloser
type TestBody struct {
	Data []byte
}

// Read implements io.Reader
func (b *TestBody) Read(p []byte) (int, error) {
	copy(p, b.Data)
	return len(p), nil
}

// Close implements io.Closer
func (b *TestBody) Close() error {
	return nil
}

// Get implements Clienter
func (c *TestClient) Get(uri string) (*http.Response, error) {
	source, ok := c.Sources[uri]
	if !ok {
		return nil, errors.New("source not found")
	}
	res := &http.Response{
		StatusCode: source.Code,
		Body:       source.Body,
	}
	<-time.After(time.Duration(source.Timeout) * time.Millisecond)
	return res, nil
}

// TestWinnerHandler testing handler for /winner GET requests
func TestWinnerHandler(t *testing.T) {
	client := &TestClient{
		Sources: map[string]*TestResponse{
			"http://example.com/fibo": &TestResponse{
				Timeout: 70,
				Code:    http.StatusOK,
				Body: &TestBody{
					Data: []byte(`[{"price":5}, {"price":3}, {"price": 8}]`),
				},
			},
			"http://example.com/primes": &TestResponse{
				Timeout: 70,
				Code:    http.StatusOK,
				Body: &TestBody{
					Data: []byte(`[{"price": 3}, {"price": 5}, {"price": 7}]`),
				},
			},
		},
	}
	req, err := http.NewRequest("GET", "/winner?s=http://example.com/fibo&s=http://example.com/primes", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	WinnerHandler(client).ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Expected status code: %d, got: %d", http.StatusOK, status)
	}

	winner := &Source{}
	err = json.NewDecoder(rec.Body).Decode(&winner)
	if err != nil {
		t.Errorf("Expected valid json, got err: %v", err)
	}

	if winner.Price != 7 || winner.URI != "http://example.com/fibo" {
		t.Errorf("Expected price: 7, source: http://example.com/fibo; got price: %d, source: %v", winner.Price, winner.URI)
	}

}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// SourceTimeout is value in milliseconds until source being ignored
const (
	SourceTimeout = 100
)

// ProcessSource gets prices from source and adds its price to list of all prices and sources
func ProcessSource(wg *sync.WaitGroup, client Clienter, uri string, sources *Sources) {
	defer wg.Done()
	prices := []*Price{}
	res, err := client.Get(uri)
	if err != nil {
		return
	}
	if res == nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		return
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&prices)
	if err != nil {
		return
	}
	for _, price := range prices {
		sources.Add(uri, price.Price)
	}
}

// WinnerHandler handles http-request with sources, gets data from it and try to find winner
func WinnerHandler(client Clienter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := r.URL.Query()["s"]
		sources := &Sources{
			Data: []*Source{},
		}
		if len(s) == 0 {
			WriteResponse(w, http.StatusBadRequest, APIResponse{
				Code:    "BadRequest",
				Message: "No one source provided",
			})
			return
		}

		wg := &sync.WaitGroup{}
		wg.Add(len(s))
		for _, uri := range s {
			go ProcessSource(wg, client, uri, sources)
		}
		wg.Wait()
		winner, err := sources.Winner()
		if err != nil {
			WriteResponse(w, http.StatusInternalServerError, APIResponse{
				Code:    "InternalServerError",
				Message: err.Error(),
			})
			return
		}
		WriteResponse(w, http.StatusOK, winner)
	})
}

// PingHandler is for health checking
func PingHandler(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, http.StatusOK, APIResponse{
		Message: "pong",
	})
}

// ForbiddenHandler handles all requests not to /winner uri
func ForbiddenHandler(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, http.StatusForbidden, APIResponse{
		Code:    "Forbidden",
		Message: "This uri is forbidden",
	})
}

func main() {
	port := flag.Int("p", 8080, "Defines TCP-port of service")
	flag.Parse()

	http.HandleFunc("/", ForbiddenHandler)
	http.HandleFunc("/ping", PingHandler)
	http.HandleFunc("/winner", WinnerHandler(&http.Client{
		Timeout: time.Duration(SourceTimeout * time.Millisecond),
	}))

	fmt.Printf("Listen on port %d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

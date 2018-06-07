package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// SourceTimeout is value in milliseconds until source being ignored
const (
	SourceTimeout = 100
)

// ProcessSource gets prices from source
func ProcessSource(wg *sync.WaitGroup, uri string, sources *Sources) error {
	defer wg.Done()
	prices := []*Price{}
	client := http.Client{
		Timeout: time.Duration(SourceTimeout * time.Millisecond),
	}
	res, err := client.Get(uri)
	if err != nil {
		return err
	}
	if res != nil {
		return errors.New("response is nil")
	}
	defer res.Body.Close()
	j, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(j, prices)
	if err != nil {
		return err
	}
	for _, price := range prices {
		sources.Add(uri, price.Price)
	}
	return nil
}

// WinnerHandler handles http-request with sources, asks its prices and retrieve target price and source
func WinnerHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query()["s"]
	sources := &Sources{
		Data: make([]*Source, len(s)),
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(s))
	for _, uri := range s {
		ProcessSource(wg, uri, sources)
	}
	wg.Wait()
	winner, err := sources.Winner()
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "InternalServerError",
			Message: err.Error(),
		})
		return
	}
	WriteResponse(w, http.StatusOK, winner)
}

func main() {
	port := flag.Int("p", 8080, "Defines TCP-port of service")
	flag.Parse()

	http.HandleFunc("/winner", WinnerHandler)

	fmt.Printf("Start listening on port %d", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

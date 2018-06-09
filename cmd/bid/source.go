package main

import (
	"errors"
	"net/http"
	"sort"
	"sync"
)

// Clienter is the interface which provide one method Get for
// retrieving json-syntah response from other sources
type Clienter interface {
	Get(uri string) (*http.Response, error)
}

// Price contains price from source
type Price struct {
	Price int `json:"price"`
}

// Source is struct for source URI and price
type Source struct {
	URI   string `json:"uri"`
	Price int    `json:"price"`
}

// Sources contains list of all sources and its prices
type Sources struct {
	sync.RWMutex
	Data []*Source
}

// Len returns length of sources slice (for sorting)
func (s *Sources) Len() int {
	return len(s.Data)
}

// Swap swaps items with i and j indexes (for sorting)
func (s *Sources) Swap(i, j int) {
	s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
}

// Less returns true if value from previous index less than next (for sorting, asc)
func (s *Sources) Less(i, j int) bool {
	return s.Data[i].Price < s.Data[j].Price
}

// Add adds source and its prices to the list
func (s *Sources) Add(uri string, price int) {
	s.Lock()
	s.Data = append(s.Data, &Source{
		URI:   uri,
		Price: price,
	})
	s.Unlock()
}

// Winner returns second highest price and source uri which wins providing highest price
func (s *Sources) Winner() (*Source, error) {
	s.Lock()
	defer s.Unlock()
	sort.Sort(s)
	count := len(s.Data)
	if count > 0 {
		uri := s.Data[count-1].URI
		price := s.Data[count-1].Price
		if count > 1 {
			price = s.Data[count-2].Price
		}
		return &Source{
			URI:   uri,
			Price: price,
		}, nil
	}
	return nil, errors.New("prices not found, there is no winner")
}

package main

import (
	"errors"
	"sort"
	"sync"
)

// Price contains price from source
type Price struct {
	Price int `json:"price"`
}

// Source is struct for source URI and price
type Source struct {
	URI   string
	Price int
}

// Sources contains list of all sources and its prices
type Sources struct {
	sync.RWMutex
	Data []*Source
}

// Len returns length of sources slice (for sorting)
func (s *Sources) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.Data)
}

// Swap swaps items with i and j indexes (for sorting)
func (s *Sources) Swap(i, j int) {
	s.Lock()
	s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
	s.Unlock()
}

// Less returns true if value from previous index less than next (for sorting, asc)
func (s *Sources) Less(i, j int) bool {
	s.RLock()
	defer s.RUnlock()
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
	s.RLock()
	defer s.RUnlock()
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

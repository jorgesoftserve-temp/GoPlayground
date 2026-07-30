package main

import "sync"

type VisitedSet struct {
	mu   sync.Mutex
	urls map[string]struct{}
}

func NewVisitedSet() *VisitedSet {
	return &VisitedSet{
		urls: make(map[string]struct{}),
	}
}

func (visited *VisitedSet) TryAdd(url string) bool {
	visited.mu.Lock()
	defer visited.mu.Unlock()

	if _, exists := visited.urls[url]; exists {
		return false
	}

	visited.urls[url] = struct{}{}

	return true
}

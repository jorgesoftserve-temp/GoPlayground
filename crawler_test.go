package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakePage struct {
	Links []string
	Err   error
}

type FakeFetcher struct {
	mu         sync.Mutex
	pages      map[string]FakePage
	fetchCount map[string]int
}

func NewFakeFetcher(pages map[string]FakePage) *FakeFetcher {
	return &FakeFetcher{
		pages:      pages,
		fetchCount: make(map[string]int),
	}
}

func (fetcher *FakeFetcher) Fetch(url string) ([]string, error) {
	fetcher.mu.Lock()
	fetcher.fetchCount[url]++
	fetcher.mu.Unlock()

	page, exists := fetcher.pages[url]
	if !exists {
		return nil, fmt.Errorf("URL no encontrada: %s", url)
	}

	return page.Links, page.Err
}

func (fetcher *FakeFetcher) Count(url string) int {
	fetcher.mu.Lock()
	defer fetcher.mu.Unlock()

	return fetcher.fetchCount[url]
}

func TestCrawlerCrawl(t *testing.T) {
	tests := []struct {
		name       string
		startURL   string
		maxDepth   int
		pages      map[string]FakePage
		expected   []string
		errorCount int
	}{
		{
			name:     "profundidad cero solo visita la URL inicial",
			startURL: "https://test.com",
			maxDepth: 0,
			pages: map[string]FakePage{
				"https://test.com": {
					Links: []string{
						"https://test.com/a",
						"https://test.com/b",
					},
				},
			},
			expected: []string{
				"https://test.com",
			},
		},
		{
			name:     "respeta la profundidad máxima",
			startURL: "https://test.com",
			maxDepth: 1,
			pages: map[string]FakePage{
				"https://test.com": {
					Links: []string{
						"https://test.com/a",
						"https://test.com/b",
					},
				},
				"https://test.com/a": {
					Links: []string{
						"https://test.com/deep",
					},
				},
				"https://test.com/b":    {},
				"https://test.com/deep": {},
			},
			expected: []string{
				"https://test.com",
				"https://test.com/a",
				"https://test.com/b",
			},
		},
		{
			name:     "no procesa URLs duplicadas",
			startURL: "https://test.com",
			maxDepth: 2,
			pages: map[string]FakePage{
				"https://test.com": {
					Links: []string{
						"https://test.com/a",
						"https://test.com/b",
					},
				},
				"https://test.com/a": {
					Links: []string{
						"https://test.com/shared",
					},
				},
				"https://test.com/b": {
					Links: []string{
						"https://test.com/shared",
					},
				},
				"https://test.com/shared": {},
			},
			expected: []string{
				"https://test.com",
				"https://test.com/a",
				"https://test.com/b",
				"https://test.com/shared",
			},
		},
		{
			name:     "continúa cuando una página devuelve error",
			startURL: "https://test.com",
			maxDepth: 1,
			pages: map[string]FakePage{
				"https://test.com": {
					Links: []string{
						"https://test.com/a",
						"https://test.com/b",
					},
				},
				"https://test.com/a": {
					Err: fmt.Errorf("falló el fetch"),
				},
				"https://test.com/b": {},
			},
			expected: []string{
				"https://test.com",
				"https://test.com/a",
				"https://test.com/b",
			},
			errorCount: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fetcher := NewFakeFetcher(test.pages)
			crawler := NewCrawler(fetcher, 5)

			report, err := crawler.Crawl(
				test.startURL,
				test.maxDepth,
			)

			require.NoError(t, err)

			assert.ElementsMatch(
				t,
				test.expected,
				report.URLs,
			)

			assert.Len(
				t,
				report.Errors,
				test.errorCount,
			)
		})
	}
}

func TestCrawlerDoesNotFetchTheSameURLTwice(t *testing.T) {
	fetcher := NewFakeFetcher(
		map[string]FakePage{
			"https://test.com": {
				Links: []string{
					"https://test.com/a",
					"https://test.com/b",
				},
			},
			"https://test.com/a": {
				Links: []string{
					"https://test.com/shared",
				},
			},
			"https://test.com/b": {
				Links: []string{
					"https://test.com/shared",
				},
			},
			"https://test.com/shared": {},
		},
	)

	crawler := NewCrawler(fetcher, 5)

	_, err := crawler.Crawl("https://test.com", 2)
	require.NoError(t, err)

	assert.Equal(
		t,
		1,
		fetcher.Count("https://test.com/shared"),
	)
}

type ConcurrencyFetcher struct {
	mu            sync.Mutex
	active        int
	maxActive     int
	numberOfLinks int
	delay         time.Duration
}

func (fetcher *ConcurrencyFetcher) Fetch(
	url string,
) ([]string, error) {
	fetcher.mu.Lock()

	fetcher.active++

	if fetcher.active > fetcher.maxActive {
		fetcher.maxActive = fetcher.active
	}

	fetcher.mu.Unlock()

	time.Sleep(fetcher.delay)

	var links []string

	if url == "https://test.com" {
		for index := 0; index < fetcher.numberOfLinks; index++ {
			links = append(
				links,
				fmt.Sprintf(
					"https://test.com/page-%d",
					index,
				),
			)
		}
	}

	fetcher.mu.Lock()
	fetcher.active--
	fetcher.mu.Unlock()

	return links, nil
}

func (fetcher *ConcurrencyFetcher) MaxActive() int {
	fetcher.mu.Lock()
	defer fetcher.mu.Unlock()

	return fetcher.maxActive
}

func TestCrawlerLimitsConcurrency(t *testing.T) {
	const workers = 5

	fetcher := &ConcurrencyFetcher{
		numberOfLinks: 20,
		delay:         20 * time.Millisecond,
	}

	crawler := NewCrawler(fetcher, workers)

	report, err := crawler.Crawl(
		"https://test.com",
		1,
	)

	require.NoError(t, err)
	require.Len(t, report.Errors, 0)

	assert.LessOrEqual(
		t,
		fetcher.MaxActive(),
		workers,
	)

	assert.Equal(
		t,
		workers,
		fetcher.MaxActive(),
	)
}

func TestCrawlerValidations(t *testing.T) {
	tests := []struct {
		name     string
		fetcher  Fetcher
		workers  int
		startURL string
		maxDepth int
	}{
		{
			name:     "fetcher nil",
			fetcher:  nil,
			workers:  5,
			startURL: "https://test.com",
			maxDepth: 1,
		},
		{
			name:     "cero workers",
			fetcher:  NewFakeFetcher(nil),
			workers:  0,
			startURL: "https://test.com",
			maxDepth: 1,
		},
		{
			name:     "URL vacía",
			fetcher:  NewFakeFetcher(nil),
			workers:  5,
			startURL: "",
			maxDepth: 1,
		},
		{
			name:     "profundidad negativa",
			fetcher:  NewFakeFetcher(nil),
			workers:  5,
			startURL: "https://test.com",
			maxDepth: -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			crawler := NewCrawler(
				test.fetcher,
				test.workers,
			)

			_, err := crawler.Crawl(
				test.startURL,
				test.maxDepth,
			)

			require.Error(t, err)
		})
	}
}

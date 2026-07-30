package main

import (
	"fmt"
	"sync"
)

type CrawlJob struct {
	URL   string
	Depth int
}

type FetchResult struct {
	Job   CrawlJob
	Links []string
	Err   error
}

type CrawlReport struct {
	URLs   []string
	Errors []error
}

type Crawler struct {
	fetcher    Fetcher
	workers    int
	newVisited func() *VisitedSet
}

func NewCrawler(fetcher Fetcher, workers int) *Crawler {
	return &Crawler{
		fetcher: fetcher,
		workers: workers,
		newVisited: func() *VisitedSet {
			return NewVisitedSet()
		},
	}
}

func (crawler *Crawler) Crawl(
	startURL string,
	maxDepth int,
) (CrawlReport, error) {
	if crawler.fetcher == nil {
		return CrawlReport{}, fmt.Errorf("el fetcher no puede ser nil")
	}

	if crawler.workers <= 0 {
		return CrawlReport{}, fmt.Errorf(
			"la cantidad de workers debe ser mayor que cero",
		)
	}

	if startURL == "" {
		return CrawlReport{}, fmt.Errorf(
			"la URL inicial no puede estar vacía",
		)
	}

	if maxDepth < 0 {
		return CrawlReport{}, fmt.Errorf(
			"la profundidad máxima no puede ser negativa",
		)
	}

	visited := crawler.newVisited()

	jobs := make(chan CrawlJob)
	results := make(chan FetchResult, crawler.workers)

	var workersWG sync.WaitGroup

	for workerID := 0; workerID < crawler.workers; workerID++ {
		workersWG.Add(1)

		go crawler.worker(
			workerID,
			jobs,
			results,
			&workersWG,
		)
	}

	startJob := CrawlJob{
		URL:   startURL,
		Depth: 0,
	}

	visited.TryAdd(startURL)

	report := CrawlReport{
		URLs: []string{startURL},
	}

	pendingJobs := []CrawlJob{startJob}
	activeJobs := 0

	for len(pendingJobs) > 0 || activeJobs > 0 {
		var jobsChannel chan CrawlJob
		var nextJob CrawlJob

		if len(pendingJobs) > 0 {
			jobsChannel = jobs
			nextJob = pendingJobs[0]
		}

		select {
		case jobsChannel <- nextJob:
			pendingJobs = pendingJobs[1:]
			activeJobs++

		case result := <-results:
			activeJobs--

			if result.Err != nil {
				report.Errors = append(
					report.Errors,
					fmt.Errorf(
						"no se pudo obtener %s: %w",
						result.Job.URL,
						result.Err,
					),
				)

				continue
			}

			if result.Job.Depth >= maxDepth {
				continue
			}

			for _, link := range result.Links {
				if !visited.TryAdd(link) {
					continue
				}

				report.URLs = append(report.URLs, link)

				pendingJobs = append(
					pendingJobs,
					CrawlJob{
						URL:   link,
						Depth: result.Job.Depth + 1,
					},
				)
			}
		}
	}

	close(jobs)
	workersWG.Wait()

	return report, nil
}

func (crawler *Crawler) worker(
	workerID int,
	jobs <-chan CrawlJob,
	results chan<- FetchResult,
	workersWG *sync.WaitGroup,
) {
	defer workersWG.Done()

	for job := range jobs {
		links, err := crawler.fetcher.Fetch(job.URL)

		results <- FetchResult{
			Job:   job,
			Links: links,
			Err:   err,
		}
	}
}

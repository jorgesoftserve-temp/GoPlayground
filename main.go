package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const maxConcurrentFetches = 5

func main() {
	if len(os.Args) != 3 {
		fmt.Printf(
			"Uso: %s <url-inicial> <profundidad-máxima>\n",
			os.Args[0],
		)
		os.Exit(1)
	}

	startURL := os.Args[1]

	maxDepth, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintln(
			os.Stderr,
			"Error: la profundidad debe ser un número entero",
		)
		os.Exit(1)
	}

	fetcher := NewHTTPFetcher(10 * time.Second)
	crawler := NewCrawler(fetcher, maxConcurrentFetches)

	report, err := crawler.Crawl(startURL, maxDepth)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Println("URLs encontradas:")

	for _, crawledURL := range report.URLs {
		fmt.Println(crawledURL)
	}

	if len(report.Errors) > 0 {
		fmt.Println("\nErrores encontrados:")

		for _, crawlError := range report.Errors {
			fmt.Fprintln(os.Stderr, "-", crawlError)
		}
	}
}

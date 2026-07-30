package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Fetcher interface {
	Fetch(rawURL string) ([]string, error)
}

type HTTPFetcher struct {
	client *http.Client
}

func NewHTTPFetcher(timeout time.Duration) *HTTPFetcher {
	return &HTTPFetcher{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (fetcher *HTTPFetcher) Fetch(rawURL string) ([]string, error) {
	response, err := fetcher.client.Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", rawURL, err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf(
			"GET %s devolvió el estado %s",
			rawURL,
			response.Status,
		)
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return []string{}, nil
	}

	links, err := extractLinks(response.Body, rawURL)
	if err != nil {
		return nil, fmt.Errorf(
			"no se pudo analizar el HTML de %s: %w",
			rawURL,
			err,
		)
	}

	return links, nil
}

func extractLinks(reader io.Reader, baseURL string) ([]string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("URL base inválida: %w", err)
	}

	document, err := html.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("HTML inválido: %w", err)
	}

	uniqueLinks := make(map[string]struct{})

	var walk func(node *html.Node)

	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attribute := range node.Attr {
				if attribute.Key != "href" {
					continue
				}

				link, ok := normalizeURL(base, attribute.Val)
				if ok {
					uniqueLinks[link] = struct{}{}
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(document)

	links := make([]string, 0, len(uniqueLinks))

	for link := range uniqueLinks {
		links = append(links, link)
	}

	return links, nil
}

func normalizeURL(base *url.URL, href string) (string, bool) {
	href = strings.TrimSpace(href)

	if href == "" {
		return "", false
	}

	reference, err := url.Parse(href)
	if err != nil {
		return "", false
	}

	resolved := base.ResolveReference(reference)

	if resolved.Scheme != "http" && resolved.Scheme != "https" {
		return "", false
	}

	resolved.Fragment = ""

	return resolved.String(), true
}

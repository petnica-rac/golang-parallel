package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const baseURL = "http://localhost:8080"

var (
	wg      sync.WaitGroup
	mu      sync.Mutex
	visited = make(map[string]bool)
)

func extractLinks(r io.Reader) []string {
	var links []string
	doc, err := html.Parse(r)
	if err != nil {
		return nil
	}
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := attr.Val
					if strings.HasPrefix(href, "/") {
						href = baseURL + href
					}
					links = append(links, href)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return links
}

func crawl(url string) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error fetching %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading %s: %v\n", url, err)
		return
	}

	links := extractLinks(bytes.NewReader(body))
	fmt.Printf("fetched %s (%d bytes, %d links)\n", url, len(body), len(links))

	for _, link := range links {
		mu.Lock()
		if !visited[link] {
			visited[link] = true
			wg.Add(1)
			go crawl(link)
		}
		mu.Unlock()
	}
}

func main() {
	seed := baseURL + "/page/1"

	start := time.Now()

	visited[seed] = true

	wg.Add(1)
	go crawl(seed)

	wg.Wait()

	fmt.Printf("\ncrawled %d pages in %v\n", len(visited), time.Since(start))
}

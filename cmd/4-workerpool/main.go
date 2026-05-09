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

const (
	baseURL    = "http://localhost:8080"
	numWorkers = 5
)

var (
	jobs    = make(chan string, 100)
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
		unseen := !visited[link]
		if unseen {
			visited[link] = true
			wg.Add(1)
		}
		mu.Unlock()

		if unseen {
			jobs <- link
		}
	}
}

func main() {
	for range numWorkers {
		go func() {
			for url := range jobs {
				crawl(url)
			}
		}()
	}

	seed := baseURL + "/page/1"

	mu.Lock()
	visited[seed] = true
	mu.Unlock()

	start := time.Now()
	wg.Add(1)
	jobs <- seed

	wg.Wait()
	close(jobs)

	fmt.Printf("\ncrawled %d pages in %v\n", len(visited), time.Since(start))
}

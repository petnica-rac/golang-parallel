package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const (
	baseURL         = "http://localhost:8080"
	numFetchWorkers = 10
)

var numCompressWorkers = runtime.NumCPU()

type page struct {
	url  string
	body []byte
}

var (
	jobs    = make(chan string, 100)
	pages   = make(chan page, 100)
	fetchWg sync.WaitGroup
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

func fetch(url string) {
	defer fetchWg.Done()

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

	fmt.Printf("fetched %s (%d bytes)\n", url, len(body))

	for _, link := range extractLinks(bytes.NewReader(body)) {
		mu.Lock()
		unseen := !visited[link]
		if unseen {
			visited[link] = true
			fetchWg.Add(1)
		}
		mu.Unlock()

		if unseen {
			jobs <- link
		}
	}

	pages <- page{url: url, body: body}
}

func compress(p page) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write(p.body)
	w.Close()
	ratio := float64(buf.Len()) / float64(len(p.body)) * 100
	fmt.Printf("compressed %s: %d → %d bytes (%.1f%%)\n", p.url, len(p.body), buf.Len(), ratio)
}

func main() {
	for range numFetchWorkers {
		go func() {
			for url := range jobs {
				fetch(url)
			}
		}()
	}

	var compressWg sync.WaitGroup
	for range numCompressWorkers {
		compressWg.Add(1)
		go func() {
			defer compressWg.Done()
			for p := range pages {
				compress(p)
			}
		}()
	}

	seed := baseURL + "/page/1"
	visited[seed] = true

	start := time.Now()
	fetchWg.Add(1)
	jobs <- seed

	fetchWg.Wait()
	close(jobs)
	close(pages)
	compressWg.Wait()

	fmt.Printf("\ncrawled %d pages in %v\n", len(visited), time.Since(start))
}

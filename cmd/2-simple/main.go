package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	baseURL      = "http://localhost:8080"
	fetchTimeout = 3 * time.Second
)

type result struct {
	url   string
	links []string
	bytes int
	err   error
}

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

func fetch(url string, ch chan<- result) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- result{url: url, err: err}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- result{url: url, err: err}
		return
	}

	ch <- result{
		url:   url,
		links: extractLinks(bytes.NewReader(body)),
		bytes: len(body),
	}
}

var urls = []string{
	"http://localhost:8080/page/1",
	"http://localhost:8080/page/2",
	"http://localhost:8080/page/3",
	"http://localhost:8080/page/4",
	"http://localhost:8080/page/5",
	"http://localhost:8080/page/6",
	"http://localhost:8080/page/7",
	"http://localhost:8080/page/8",
	"http://localhost:8080/page/9",
	"http://localhost:8080/page/10",
	"http://localhost:8080/page/11",
	"http://localhost:8080/page/12",
	"http://localhost:8080/page/13",
	"http://localhost:8080/page/14",
	"http://localhost:8080/page/15",
	"http://localhost:8080/page/16",
	"http://localhost:8080/page/17",
	"http://localhost:8080/page/18",
	"http://localhost:8080/page/19",
	"http://localhost:8080/page/20",
}

func main() {
	ch := make(chan result)

	start := time.Now()
	for _, url := range urls {
		go fetch(url, ch)
	}

	for range urls {
		select {
		case r := <-ch:
			if r.err != nil {
				fmt.Printf("error fetching %s: %v\n", r.url, r.err)
				continue
			}
			fmt.Printf("fetched %s (%d bytes, %d links)\n", r.url, r.bytes, len(r.links))
			for _, link := range r.links {
				fmt.Printf("  -> %s\n", link)
			}
		case <-time.After(fetchTimeout):
			fmt.Println("timeout: a fetch took too long, moving on")
		}
	}

	fmt.Printf("\ndone in %v\n", time.Since(start))
}

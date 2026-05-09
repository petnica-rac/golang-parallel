package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	port      = flag.Int("port", 8080, "port to listen on")
	minDelay  = flag.Duration("min-delay", 200*time.Millisecond, "minimum response delay")
	maxDelay  = flag.Duration("max-delay", 500*time.Millisecond, "maximum response delay")
	minSize   = flag.Int("min-size", 512, "minimum response body size in bytes")
	maxSize   = flag.Int("max-size", 2048, "maximum response body size in bytes")
	pageCount = flag.Int("pages", 20, "total number of pages to serve")
	slowProb  = flag.Float64("slow-prob", 0, "probability (0.0-1.0) that a request also sleeps for slow-delay")
	slowDelay = flag.Duration("slow-delay", 5*time.Second, "extra delay applied to slow requests")
)

var pageTmpl = template.Must(template.New("page").Parse(`<!DOCTYPE html>
<html>
<head><title>Page {{.N}}</title></head>
<body>
<h1>Page {{.N}}</h1>
<p>{{.Content}}</p>
<ul>
{{range .Links}}<li><a href="/page/{{.}}">Page {{.}}</a></li>
{{end}}</ul>
</body>
</html>
`))

type pageData struct {
	N       int
	Content string
	Links   []int
}

// linksForPage returns a deterministic list of links for page n so that
// the graph is stable across server restarts and crawl runs.
func linksForPage(n, total int) []int {
	rng := rand.New(rand.NewSource(int64(n)))
	count := 3 + rng.Intn(3) // 3 to 5 links per page
	seen := map[int]bool{n: true}
	var links []int
	for len(links) < count {
		p := 1 + rng.Intn(total)
		if !seen[p] {
			seen[p] = true
			links = append(links, p)
		}
	}
	return links
}

func randDelay() {
	delta := int64(*maxDelay - *minDelay)
	d := *minDelay
	if delta > 0 {
		d += time.Duration(rand.Int63n(delta))
	}
	time.Sleep(d)
}

func handlePage(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(r.PathValue("n"))
	if err != nil || n < 1 || n > *pageCount {
		http.NotFound(w, r)
		return
	}

	randDelay()
	if rand.Float64() < *slowProb {
		time.Sleep(*slowDelay)
	}

	size := *minSize
	if delta := *maxSize - *minSize; delta > 0 {
		size += rand.Intn(delta)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pageTmpl.Execute(w, pageData{
		N:       n,
		Content: strings.Repeat("x", size),
		Links:   linksForPage(n, *pageCount),
	})
}

func main() {
	flag.Parse()

	if *minDelay > *maxDelay {
		log.Fatal("min-delay must be <= max-delay")
	}
	if *minSize > *maxSize {
		log.Fatal("min-size must be <= max-size")
	}
	if *slowProb < 0 || *slowProb > 1 {
		log.Fatal("slow-prob must be between 0.0 and 1.0")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/page/{n}", handlePage)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("listening on %s | delay %v-%v | body %d-%d bytes | pages %d | slow-prob %.2f slow-delay %v",
		addr, *minDelay, *maxDelay, *minSize, *maxSize, *pageCount, *slowProb, *slowDelay)
	log.Fatal(http.ListenAndServe(addr, mux))
}

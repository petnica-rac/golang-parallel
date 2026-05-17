package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
		log.Fatalf("invalid value for %s: %q", key, v)
	}
	return def
}

func envFloat64(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
		log.Fatalf("invalid value for %s: %q", key, v)
	}
	return def
}

func envDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
		log.Fatalf("invalid value for %s: %q", key, v)
	}
	return def
}

var (
	port      = flag.Int("port", envInt("PORT", 8080), "port to listen on")
	minDelay  = flag.Duration("min-delay", envDuration("MIN_DELAY", 200*time.Millisecond), "minimum response delay")
	maxDelay  = flag.Duration("max-delay", envDuration("MAX_DELAY", 500*time.Millisecond), "maximum response delay")
	minSize   = flag.Int("min-size", envInt("MIN_SIZE", 512), "minimum response body size in bytes")
	maxSize   = flag.Int("max-size", envInt("MAX_SIZE", 2048), "maximum response body size in bytes")
	pageCount = flag.Int("pages", envInt("PAGES", 20), "total number of pages to serve")
	slowProb  = flag.Float64("slow-prob", envFloat64("SLOW_PROB", 0), "probability (0.0-1.0) that a request also sleeps for slow-delay")
	slowDelay = flag.Duration("slow-delay", envDuration("SLOW_DELAY", 5*time.Second), "extra delay applied to slow requests")
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

// wordList is a set of common English words used to generate realistic page content.
var wordList = []string{
	"the", "be", "to", "of", "and", "a", "in", "that", "have", "it",
	"for", "not", "on", "with", "as", "you", "do", "at", "this", "but",
	"by", "from", "they", "we", "say", "she", "or", "an", "will", "my",
	"one", "all", "would", "there", "their", "what", "so", "up", "out",
	"if", "about", "who", "get", "which", "go", "when", "make", "can",
	"like", "time", "no", "just", "know", "take", "people", "into", "year",
	"your", "good", "some", "could", "them", "see", "other", "than", "then",
	"now", "look", "only", "come", "over", "think", "also", "back", "after",
	"use", "two", "how", "our", "work", "first", "well", "way", "even",
	"new", "want", "because", "any", "these", "give", "day", "most", "us",
	"great", "between", "need", "large", "often", "hand", "high", "place",
	"hold", "turn", "where", "much", "before", "move", "right", "boy",
	"old", "too", "same", "tell", "does", "set", "three", "want", "air",
	"play", "small", "number", "off", "always", "next", "show", "every",
	"near", "add", "food", "between", "own", "below", "country", "plant",
	"last", "school", "father", "keep", "tree", "never", "start", "city",
	"earth", "eye", "light", "thought", "head", "under", "story", "saw",
	"left", "don't", "few", "while", "along", "might", "close", "something",
	"seem", "next", "hard", "open", "example", "begin", "life", "always",
}

func randText(size int) string {
	var b strings.Builder
	for b.Len() < size {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(wordList[rand.Intn(len(wordList))])
	}
	return b.String()[:size]
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
		Content: randText(size),
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

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	port     = flag.Int("port", 8080, "port to listen on")
	minDelay = flag.Duration("min-delay", 200*time.Millisecond, "minimum response delay")
	maxDelay = flag.Duration("max-delay", 500*time.Millisecond, "maximum response delay")
	minSize  = flag.Int("min-size", 512, "minimum response body size in bytes")
	maxSize  = flag.Int("max-size", 2048, "maximum response body size in bytes")
)

func randDelay() {
	delta := int64(*maxDelay - *minDelay)
	d := *minDelay
	if delta > 0 {
		d += time.Duration(rand.Int63n(delta))
	}
	time.Sleep(d)
}

func randBody() string {
	size := *minSize
	if delta := *maxSize - *minSize; delta > 0 {
		size += rand.Intn(delta)
	}
	return strings.Repeat("x", size)
}

func handlePage(w http.ResponseWriter, r *http.Request) {
	randDelay()
	fmt.Fprintln(w, randBody())
}

func main() {
	flag.Parse()

	if *minDelay > *maxDelay {
		log.Fatal("min-delay must be <= max-delay")
	}
	if *minSize > *maxSize {
		log.Fatal("min-size must be <= max-size")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/page/{n}", handlePage)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("listening on %s | delay %v-%v | body %d-%d bytes",
		addr, *minDelay, *maxDelay, *minSize, *maxSize)
	log.Fatal(http.ListenAndServe(addr, mux))
}

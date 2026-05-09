package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func fetch(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error fetching %s: %v\n", url, err)
		return
	}

	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("fetched %s (%d bytes).\n", url, len(body))

	resp.Body.Close()
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
	start := time.Now()
	for _, url := range urls {
		fetch(url)
	}
	time.Sleep(time.Second * 10)
	fmt.Printf("\ndone in %v\n", time.Since(start))
}

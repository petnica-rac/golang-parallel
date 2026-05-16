package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

const dir = "texts"
const numWorkers = 4

func main() {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "readdir: %v\n", err)
		os.Exit(1)
	}

	jobs := make(chan string, len(entries))
	results := make(chan map[string]int, len(entries))

	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				results <- countWords(path)
			}
		}()
	}

	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".txt") {
			jobs <- filepath.Join(dir, e.Name())
		}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	freq := make(map[string]int)
	for m := range results {
		for w, c := range m {
			freq[w] += c
		}
	}

	type entry struct {
		word  string
		count int
	}
	var sorted []entry
	for w, c := range freq {
		sorted = append(sorted, entry{w, c})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].count > sorted[j].count
	})

	for _, e := range sorted[:min(20, len(sorted))] {
		fmt.Printf("%-15s %d\n", e.word, e.count)
	}
}

func countWords(path string) map[string]int {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	counts := make(map[string]int)
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		counts[scanner.Text()]++
	}
	return counts
}

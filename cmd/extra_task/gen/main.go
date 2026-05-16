package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

var words = []string{
	"the", "a", "and", "of", "to", "in", "is", "that", "it", "was",
	"for", "on", "are", "with", "as", "at", "be", "this", "have", "from",
	"or", "one", "had", "by", "but", "not", "what", "all", "were", "when",
	"we", "there", "can", "an", "your", "which", "their", "if", "do", "will",
	"each", "about", "how", "up", "out", "them", "then", "she", "many", "some",
	"so", "these", "would", "other", "into", "has", "more", "her", "two", "like",
	"him", "see", "time", "could", "no", "make", "than", "first", "been", "its",
	"who", "now", "people", "my", "made", "over", "did", "down", "only", "way",
	"find", "use", "may", "water", "long", "little", "very", "after", "word", "called",
	"just", "where", "most", "know", "get", "through", "back", "much", "before", "go",
	"good", "new", "write", "our", "used", "me", "man", "too", "any", "day",
	"same", "right", "look", "think", "also", "around", "came", "three", "small", "set",
	"put", "end", "does", "another", "well", "large", "need", "big", "high", "such",
	"follow", "act", "why", "ask", "men", "change", "went", "light", "kind", "off",
	"play", "spell", "air", "away", "animal", "house", "point", "page", "letter", "mother",
}

func main() {
	outDir := flag.String("out", "texts", "output directory")
	nFiles := flag.Int("files", 20, "number of files to generate")
	nWords := flag.Int("words", 2000, "words per file")
	flag.Parse()

	if err := os.MkdirAll(*outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
		os.Exit(1)
	}

	for i := range *nFiles {
		path := filepath.Join(*outDir, fmt.Sprintf("text%02d.txt", i+1))
		if err := writeFile(path, *nWords); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", path, err)
			os.Exit(1)
		}
	}

	fmt.Printf("wrote %d files to %s\n", *nFiles, *outDir)
}

func writeFile(path string, n int) error {
	buf := make([]string, n)
	for i := range n {
		buf[i] = words[rand.Intn(len(words))]
	}
	return os.WriteFile(path, []byte(strings.Join(buf, " ")), 0644)
}

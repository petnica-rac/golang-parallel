# Parallel Programming in Go

Workshop materials for a 6-hour hands-on introduction to concurrent programming in Go, developed for Petnica Science Center 2026.

## Overview

Students build a concurrent web crawler from scratch across six progressive blocks. Each block introduces a new concept and adds a layer to the same program — nothing gets thrown away. By the end, the crawler fetches pages in parallel, follows links, limits concurrency, compresses content in a pipeline, and shuts down gracefully on timeout or interrupt.

A local test server (included) simulates realistic network latency and serves a graph of HTML pages for the crawler to explore.

## Concepts covered

- **Goroutines** — launching concurrent work and understanding why sequential I/O is wasteful
- **Channels** — communicating between goroutines, collecting results, directional channel types
- **Select** — handling multiple channel operations and timeouts
- **sync.WaitGroup** — waiting for a dynamic number of goroutines to finish
- **sync.Mutex** — protecting shared state from concurrent access
- **Worker pools** — capping concurrency with a fixed number of long-lived goroutines
- **Pipelines** — chaining independent processing stages with channels
- **context.Context** — propagating cancellation and timeouts through a call tree
- **Race detection** — using `-race` to surface data races that don't produce obvious wrong output

## Structure

```
cmd/
  server/       — test server (run this before any crawler)
  1-naive/      — sequential fetcher, then goroutines with time.Sleep
  2-simple/     — channels, select, HTML link parsing
  3-recursive/  — WaitGroup, Mutex, recursive crawl
  4-workerpool/ — fixed worker pool
  5-pipeline/   — two-stage pipeline with gzip compression
  6-context/    — timeout and graceful cancellation
doc/
  plan.md       — instructor guide and block-by-block workshop plan
  server.md     — server configuration reference
  materials/    — per-block reference sheets given to students
```

## Running

```
# start the test server
go run ./cmd/server

# run any crawler
go run ./cmd/4-workerpool
```

See `doc/server.md` for server configuration options.

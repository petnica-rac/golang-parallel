# Block 4 — Reference Sheet

Replace the unbounded goroutine-per-URL approach from Block 3 with a fixed pool of worker goroutines. Instead of launching a new goroutine for every discovered URL, send URLs into a shared channel and have a fixed number of workers read from it. This keeps the total number of goroutines constant regardless of how many pages are crawled.

## Buffered channels

```go
ch := make(chan string, 100) // can hold up to 100 items without a receiver ready

ch <- "value"  // does not block if buffer has space
v := <-ch      // does not block if buffer has items
```

## Worker pool pattern

```go
jobs := make(chan string, 100)

// start N workers
for range N {
    go func() {
        for job := range jobs {
            process(job)
        }
    }()
}

// send work
jobs <- "url1"
jobs <- "url2"

// signal workers to stop
close(jobs)
```

## Closing a channel

```go
close(ch)
```

- Sending to a closed channel panics
- Receiving from a closed, empty channel returns the zero value and `false`
- `for range ch` exits cleanly when the channel is closed and drained

## Command-line flags (stretch)

```go
import "flag"

n := flag.Int("workers", 5, "number of worker goroutines")
flag.Parse()

fmt.Println(*n) // dereference the pointer
```

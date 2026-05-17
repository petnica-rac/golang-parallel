# Test Server

The workshop uses a local Go HTTP server that all students point their crawlers at. It serves a graph of HTML pages with internal links, simulates realistic network latency, and can be configured to misbehave in ways that motivate specific concurrency solutions.

Run locally with:

```
go run ./cmd/server [flags]
```

## Deployment

The server is deployed on fly.io. To redeploy after changing configuration:

```
fly deploy
```

Environment variables can be set in the `[env]` section of `fly.toml` — all flags have a corresponding env var (same name, uppercase). For example:

```toml
[env]
  SLOW_PROB = "0.2"
  SLOW_DELAY = "5s"
  MIN_DELAY = "1s"
  MAX_DELAY = "2s"
  PAGES = "30"
```

Update `fly.toml` to match the desired configuration for each block, then run `fly deploy` before the session.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `8080` | Port to listen on |
| `-pages` | `20` | Total number of pages in the graph |
| `-min-delay` | `200ms` | Minimum artificial delay per request |
| `-max-delay` | `500ms` | Maximum artificial delay per request |
| `-min-size` | `512` | Minimum response body size in bytes |
| `-max-size` | `2048` | Maximum response body size in bytes |
| `-slow-prob` | `0.0` | Probability (0.0–1.0) that a request also sleeps for `-slow-delay` |
| `-slow-delay` | `5s` | Extra delay applied to slow requests |

## Endpoints

`GET /page/{n}` — serves page `n` (1-indexed, up to `-pages`). Returns 404 for out-of-range values.

Each page is a valid HTML document containing:
- A `<p>` tag with randomly generated English text (size controlled by `-min-size` / `-max-size`)
- A `<ul>` of 3–5 `<a href>` links to other pages in the graph

The link graph is deterministic — page `n` always links to the same set of pages regardless of server restarts or request order. This keeps crawl behaviour reproducible across students.

## Suggested configurations per block

**Block 1** — students are fetching a fixed list sequentially, then with goroutines. The contrast between sequential and parallel should be dramatic.
```
go run ./cmd/server -slow-prob 0
```

**Block 2** — introduce slow responses to motivate the `select` timeout exercise.
```
go run ./cmd/server -slow-prob 0.2 -slow-delay 5s
```

**Block 3 onwards** — default settings are generally fine. Adjust page count upward once students are crawling the graph dynamically.

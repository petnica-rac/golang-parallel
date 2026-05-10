# Parallel Programming in Go — Workshop Plan

**Duration:** 6 hours (lunch break after hour 4)
**Audience:** High school students, strong programming background, introductory Go knowledge
**Format:** Hands-on. Students write code throughout. Brief conceptual framing before each section, then into the editor.

## The Thread

Students are not working through isolated exercises. They are building a single program — a concurrent web crawler / link checker — that grows more capable with every block. Each block adds a layer. Nothing gets thrown away.

The crawler targets a local test server (provided) that simulates realistic network latency. This makes the performance difference between sequential and parallel execution tangible and dramatic.

---

## Block 1 — Why Concurrency? (~45 min)

### Goal

Students should finish this block with a clear, felt understanding of why concurrency matters — not because someone told them, but because they watched their program be slow and then made it fast. They should also have their first goroutine running, and a genuine open question in their heads: *"this works, but something feels wrong."*

### Flow

**1. The slow program**
Students are given a working sequential fetcher. It fetches a fixed list of URLs one by one and prints how long it took. They run it. It's slow. The instructor asks: why is most of that time being wasted?

**2. The concept**
Brief explanation of goroutines — what they are, how they relate to OS threads, what `go` does. Keep it short. The goal is just enough to make sense of the next step.

**3. The core assignment**
Students add `go` to the fetch call and run the program again. It finishes instantly — but prints nothing. The instructor leads a discussion: where did the output go? Why did the program exit? Students then add a sleep at the end of main to keep the program alive long enough to see results. It works. It's fast. But it feels fragile.

The session ends here, deliberately. Students have a working solution that they should feel uneasy about. The right question — *"what if the fetches take longer than my sleep?"* — is the bridge into Block 2.

### Stretch Goals

- Add timestamps to each fetch result to visualize goroutines running in parallel rather than in sequence.
- Introduce a counter that tracks how many URLs were successfully fetched. The count may look correct — but run the program with `-race` and observe what Go has to say about it. (This foreshadows Block 3.)

---

## Block 2 — Channels & Select (~75 min)

### Goal

Students replace the `time.Sleep` hack with a proper communication mechanism. By the end of this block they have a crawler that correctly waits for all goroutines, handles slow or unresponsive URLs gracefully, and can discover new URLs on its own by parsing the pages it fetches.

### Flow

**1. The concept**
Brief framing: goroutines need a way to talk back to the caller. Channels are Go's answer — typed, synchronized pipes. Cover the basics: creating a channel, sending, receiving. Keep it short; students will learn the rest by doing.

**2. Core assignment — replacing the sleep**
Students replace the `time.Sleep` with a results channel. Each goroutine sends its result on the channel; main receives exactly as many times as there are goroutines. The program now correctly waits for all fetches without any guesswork about timing.

**3. Core assignment — handling slow responses with `select`**
The server is configured to occasionally hang a request for much longer than normal. Students observe that the crawler can get stuck. Introduce `select` with `time.After` as a per-fetch timeout: if a fetch doesn't respond in time, report it as timed out and move on. This is the natural motivator for `select` — students feel the problem before seeing the solution.

**4. Guided — parsing links from responses**
The server now returns real HTML pages with links to other pages. The instructor walks everyone through adding a link parser using `golang.org/x/net/html`. This is not a concurrency exercise — it is about making the crawler real. Everyone implements this step together. By the end, the crawler is discovering URLs dynamically rather than working from a hardcoded list.

### Stretch Goals

- Experiment with buffered vs unbuffered channels — what changes? When does a send block?
- Change the result type from a plain string to a struct carrying the URL, body, and any error — start shaping the type the crawler will use for the rest of the workshop.

---

## Block 3 — WaitGroup & Mutex (~60 min)

### Goal

Students fix both open problems from Block 2 — the hardcoded goroutine count and the lack of cycle prevention — and end up with a crawler that genuinely drives itself by following links. This is the first version of the program that feels like a real crawler rather than a parallel fetcher.

> **Note for instructor:** This block has more moving parts than any previous one. The shift from a flat fetch loop to a recursive concurrent crawl is a meaningful conceptual jump. Prepare to spend time walking students through the overall structure before they start writing. The concepts themselves (WaitGroup, Mutex) are simple — the challenge is how they fit together in a recursive pattern.

### Flow

**1. The concept**
WaitGroup: Add before launching a goroutine, Done when it finishes, Wait blocks until the counter reaches zero. Mutex: Lock before touching shared state, Unlock when done. Keep the explanation brief — both primitives are intuitive once students see them in context.

**2. Core assignment — WaitGroup**
Replace the `for range urls` counting loop with a WaitGroup. This is a clean mechanical change. Students should feel how much more natural it is compared to counting goroutines manually or sleeping.

**3. Core assignment — Mutex and recursive crawl**
Add a visited map protected by a Mutex. Start from a single seed URL, follow discovered links, and launch a new goroutine for each unvisited URL. The crawler now drives itself. Students should observe it working — and then notice the goroutine count climbing without bound as the graph is explored. That observation is the bridge into Block 4.

### Stretch Goals

- Print `runtime.NumGoroutine()` periodically to watch goroutine count grow during a crawl.
- Swap `sync.Mutex` for `sync.RWMutex` on the visited map — when does the distinction matter?
- Add a depth limit to prevent the crawler from going too deep into the graph.

---

## Block 4 — Worker Pool (~60 min)

### Goal

Students fix the unbounded goroutine problem from Block 3. The crawler stops spawning a goroutine per URL and instead feeds work into a fixed pool of workers. By the end, concurrency is explicit and controlled — and students have a complete, well-behaved crawler before lunch.

### Flow

**1. The concept**
Brief framing: spawning one goroutine per unit of work is fine for small inputs, but doesn't scale. A worker pool caps concurrency at a fixed number regardless of input size. The jobs channel becomes the queue; workers are long-lived goroutines that pull from it.

**2. Core assignment**
Replace the per-URL goroutine launch with a fixed pool of worker goroutines reading from a buffered jobs channel. The main challenge students will encounter is termination: how does the program know when to close the jobs channel? Walk them through tracking pending jobs with a WaitGroup (separate from the worker lifecycle), closing the channel in a goroutine so main can still block cleanly.

**3. Wrap-up**
By the end of this block students have a crawler that is concurrent, safe, and resource-controlled. A good moment to step back and look at the whole program — it has come a long way from the sequential fetcher in Block 1.

### Stretch Goals

- Make `numWorkers` a command-line flag. Experiment with different values — does more workers always mean faster? Where does the server become the bottleneck?
- Set `numWorkers` to 1 and time the crawl. Compare to Block 1's sequential run. Why is it not identical?

---

**— Lunch Break —**

---

## Block 5 — Pipeline (~60 min)

### Goal

Students extend the crawler with a second processing stage, forming a pipeline. The key insight is that different stages have different concurrency characteristics — IO-bound and CPU-bound work call for different worker counts and for different reasons. By the end, students have a crawler that not only fetches pages but does meaningful work on their content.

### Flow

**1. The concept**
Brief framing: a pipeline connects stages with channels. Each stage runs independently at its own pace. Data flows through — one stage's output is the next stage's input. Contrast with the worker pool: a pool parallelises one task, a pipeline parallelises a sequence of tasks.

**2. Core assignment**
Split the fetch worker's responsibility: instead of doing everything inline, fetch workers now send `{url, body}` pairs on a channel to a second pool of compression workers. Compression workers gzip the body and record the compression ratio. Each pool has its own worker count constant.

**3. The discussion**
Once it works, ask students: what's the right number of workers for each stage? Fetch workers spend most of their time blocked on the network — more workers is fine, the CPU sits idle while they wait. Compress workers are pure CPU — beyond `runtime.NumCPU()` workers you gain nothing and start paying context switching overhead. This is where `GOMAXPROCS` becomes a concrete, observable concept rather than an abstract one.

### Stretch Goals

- Experiment with compress worker counts — time the crawl at 1, 2, 4, and `runtime.NumCPU()` workers. Plot or note where adding more stops helping.
- Set `GOMAXPROCS=1` and observe the effect on the compression phase specifically.

---

## Block 6 — Context & Cancellation (~45 min)

### Goal

Students add graceful cancellation to the crawler. The program currently has no way to stop cleanly — it either finishes or gets killed. By the end of this block it can time itself out automatically and respond to Ctrl+C with a clean shutdown. This is also a natural return to `select`, now used not as a timeout hack but as the core mechanism for watching both work and cancellation at the same time.

### Flow

**1. The concept**
Brief framing: `context.Context` is Go's standard way to propagate a cancellation signal through a call tree. Any function that does IO should accept and respect a context — this is a pattern students will use in every real Go program they write.

**2. Core assignment — timeout**
Wrap the crawl in a `context.WithTimeout`. Pass the context through to HTTP requests via `http.NewRequestWithContext` so in-flight fetches are actually cancelled when the deadline fires, not just abandoned. Update the worker loops to `select` on both the jobs channel and `ctx.Done()` so they stop picking up new work the moment cancellation fires.

**3. Core assignment — signal handling**
Wire up `signal.NotifyContext` so Ctrl+C triggers the same cancellation path. The program now shuts down cleanly on interrupt rather than dying abruptly. Students should observe the difference: the summary line still prints, in-flight requests are cancelled, and no goroutines are left hanging.

**4. Wrap-up**
Step back and look at the full program. It started as a sequential fetcher that took ten seconds to fetch twenty URLs. It is now a concurrent, rate-limited, pipelined crawler with graceful cancellation. Each block added one idea, and none of it was thrown away.

### Stretch Goals

- `context.WithDeadline` vs `context.WithTimeout` — what is the practical difference?
- **(End-of-workshop)** Add a word frequency stage aside compression — find the top 10 most common words across all crawled pages.

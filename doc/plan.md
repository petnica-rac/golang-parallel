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

*To be planned.*

---

## Block 3 — WaitGroup & Mutex (~60 min)

*To be planned.*

---

## Block 4 — Worker Pool (~60 min)

*To be planned.*

---

**— Lunch Break —**

---

## Block 5 — Pipeline (~60 min)

*To be planned.*

---

## Block 6 — Context & Cancellation (~45 min)

*To be planned.*

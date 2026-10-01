# loganalyzer

A command-line log analysis tool written in Go. It streams one or more log files, filters by log level, and reports statistics on matching endpoints and levels, without loading entire files into memory.

This is **Project 0** of a self-directed 12-month software engineering roadmap (September 2026 – August 2027), focused on building Big-Tech-level engineering readiness across backend, systems, and distributed-systems fundamentals. September's theme is programming mastery, Linux, and Git, and this tool was the main build for that month.

## Features

- **Multi-file support** — accepts one or more comma-separated log files in a single run
- **Level filtering** — show only lines matching a given log level (e.g. `ERROR`)
- **Statistics** — reports counts per endpoint and per level across all matched lines
- **Streaming** — reads line by line rather than loading a file into memory, verified against a 10GB, 182-million-line file with flat memory usage
- **Concurrent multi-file processing** — each file is read and parsed on its own goroutine, with results merged safely into shared totals
- **Malformed-line handling** — empty, truncated, or irregularly-spaced lines are skipped rather than crashing the program
- **Clear error feedback** — missing file flag and unmatched level filters produce specific, actionable messages rather than silent failure

## Usage

```
go run main.go -files=<path1,path2,...> [-level=LEVEL]
```

Examples:

```
go run main.go -files=server.log
go run main.go -files=server1.log,server2.log
go run main.go -files=server.log -level=ERROR
```

Expected log line format (whitespace-separated, at least 4 fields):

```
<timestamp> <LEVEL> <endpoint> <status>
```

## Design decisions

**`bufio.Scanner` with a raised buffer, over `bufio.Reader`.**
Go's `bufio.Scanner` streams a file line by line, holding only one line in memory at a time, which already satisfies the project's streaming requirement. Its default max line length is 64KB, which would fail on an unusually long single line (e.g. an embedded stack trace). Rather than switching to the uncapped `bufio.Reader`, the scanner's buffer was raised to 1MB via `scanner.Buffer()`, a smaller change that comfortably covers realistic log lines while keeping the simpler `Scan()`/`Text()` API. A line over 1MB would still fail; `bufio.Reader` remains the documented path forward if that becomes a real requirement.

**Goroutines per file, not per line.**
Benchmarking showed that parsing a single line (`parseLine`) takes ~124.5 ns/op in isolation, while the full 10GB / 182-million-line run took ~4 minutes end to end. At ~124.5ns per call, 182 million calls to `parseLine` alone account for roughly 22 seconds, about 9% of total runtime. The remaining ~91% is disk I/O, which is inherently sequential within a single file. Parallelizing line-level work would therefore add complexity and concurrency risk (shared maps, locking) without addressing the actual bottleneck. Concurrency is instead applied at the file level: each file is independent, so reading multiple files simultaneously (via goroutines) genuinely overlaps I/O that would otherwise happen one file at a time.

**`sync.Mutex` over channels for merging results.**
Each goroutine builds its own local counts (no shared state during a file's read) and only needs to fold those totals into shared combined maps once, at the end. A mutex around that single merge step was chosen over a channel-based collector pattern for simplicity, given there's exactly one shared resource and one short critical section. Verified race-free with `go test -race`.

## Concepts introduced and used

- CLI argument parsing (`flag`)
- File I/O and error handling as return values, not exceptions (`os.Open`, `err != nil`)
- Buffered, streaming file reads (`bufio.Scanner`, `scanner.Buffer`)
- String parsing and whitespace-safe splitting (`strings.Fields` vs `strings.Split`)
- Structs as typed records (`LogEntry`)
- Maps as frequency counters, including the zero-value idiom (`counts[key]++` with no existence check)
- Multiple return values and the `(value, ok)` pattern for fallible operations
- Goroutines, `sync.WaitGroup` for coordinating concurrent completion, and `sync.Mutex` for protecting shared state
- The Go race detector (`go test -race`) as a verification tool, not just a correctness assumption
- Table-driven unit tests (`go test`)
- Microbenchmarking (`go test -bench`) and using benchmark data to validate a design decision rather than guess at it

## Benchmarks

Measured on an Intel Core i7-4900MQ (WSL2 / Ubuntu):

| Test                                            | Result                                                                          |
| ----------------------------------------------- | ------------------------------------------------------------------------------- |
| `parseLine` (isolated)                          | ~124.5 ns/op (`go test -bench=.`)                                               |
| Full run: 10GB file, 182M lines, `-level=ERROR` | ~4 minutes, memory stayed flat (verified via `top`), no swap usage              |
| Parsing's share of total runtime                | ~9% (22s of ~240s) — confirms the bottleneck is disk I/O, not CPU-bound parsing |

## Known limitations

- Lines longer than 1MB will still fail (`bufio.ErrTooLong`); see design notes above
- No feedback distinguishes "level never appeared" from "level flag typo'd"; both produce the same message
- Field positions are fixed by index (timestamp, level, endpoint, status); no support for alternate log formats

# Gehenna
A CSV processing engine built in Go, focused on performance, reliability, and accounting for every record.
<img width="900" height="500" alt="image" src="https://github.com/user-attachments/assets/b52eb48b-8f86-433a-a496-f4a6dbd859e1" />

## Why this project?

At Track.co, processing large CSV files containing customer data for CX survey distribution was a significant engineering challenge.

Gehenna revisits that problem as a hands-on exploration of Go performance engineering: establishing baselines, writing benchmarks, interpreting pprof profiles, and measuring the effects of each optimization.

## Goals

- Process files larger than available memory.
- Validate customer data and report actionable errors.
- Generate reports of accepted, rejected, and duplicate records.
- Keep memory usage predictable as file sizes grow.
- Recover from interruptions without losing progress.
- Support retries without duplicating downstream actions.

These are project goals, not guarantees of the current implementation.

## Input

CSV files contain the following columns:

| Column | Description |
| --- | --- |
| customer_id | Customer identifier, preserved as a string |
| name | Customer name |
| email | Email address |
| phone | Phone number, preserved as a string |
| locale | Customer locale |
| campaign_id | Survey campaign identifier |
| notes | Additional customer information |

Input may contain Unicode characters, quoted commas, escaped quotes, and embedded newlines.

## Processing outcomes

Every input record should have an explicit outcome:

- **Accepted:** passes validation and duplicate checks.
- **Rejected:** fails validation, with the reason recorded.
- **Duplicate:** matches a previously accepted customer identifier.

Acceptance means the record is eligible for further processing. It does not confirm survey delivery.

## Development datasets

| File | Data records | Approximate size |
| --- | ---: | ---: |
| customers_dev.csv | 10,000 | 1.76 MB |
| customers_perf.csv | 5,000,000 | 896 MB |

Both datasets contain synthetic customer information, deliberately invalid fields, and duplicate records. They also exercise CSV quoting and multiline fields.

Use a mock service for delivery experiments.

See `DATASET_README.org` for validation rules, expected outcome counts, and profiling notes.

## Roadmap

- [ ] Sequential streaming CSV reader
- [ ] Customer validation
- [ ] Accepted, rejected, and duplicate reports
- [ ] Correctness tests against known dataset outcomes
- [ ] Baseline benchmarks and profiling
- [ ] Optimizations supported by measurements
- [ ] Mock survey distribution service
- [ ] Bounded concurrency and backpressure
- [ ] Timeouts and bounded retries
- [ ] Persistent progress and crash recovery
- [ ] Idempotent downstream processing

## Performance engineering

Begin with a correct sequential implementation. Measure before changing it.

Benchmark parsing, validation, and the complete import separately. Track:

- Records processed per second
- Input throughput in MB/s
- Allocations and allocated bytes per operation
- Peak memory usage
- End-to-end processing time

Keep fixture generation outside the timed workload. Document whether a benchmark includes disk I/O, report writing, deduplication, or downstream requests.

### Run benchmarks

After adding benchmark functions:

    go test ./... -run '^$' -bench . -benchmem -count=5

### Capture CPU and memory profiles

Run from the package containing `BenchmarkImport`:

    go test . -run '^$' -bench '^BenchmarkImport$' -benchtime=1x \
      -cpuprofile cpu.out -memprofile mem.out

Inspect the profiles:

    go tool pprof -http=localhost:8080 cpu.out
    go tool pprof -http=localhost:8081 -alloc_space mem.out

### Check for data races

Run separately from performance measurements:

    go test -race ./...

## Reliability experiments

Introduce failures deliberately:

- Stop the process during an import.
- Make the mock destination slow or unavailable.
- Return temporary errors from the mock destination.
- Resume an interrupted job.
- Verify that every record is accounted for.
- Verify that retries and recovery do not duplicate mock deliveries.

## Guiding principle

Every record has an outcome. Every optimization needs evidence.

# Database Benchmark: TimescaleDB vs ClickHouse

- **TimescaleDB Benchmarks:**
  - Single insert performance
  - Batch insert with transactions

- **ClickHouse Benchmarks:**
  - Single insert performance
  - Async insert performance (non-blocking)
  - Batch insert performance
  - Async batch insert performance

## Bench

```sh
$ make start
$ make benchmark
$ make clean
```

## Result

```sh
$ make benchmar
go test -v -bench=. -benchmem -benchtime=10s
=== RUN   TestDatabaseConnections
=== RUN   TestDatabaseConnections/TimescaleDB
    benchmark_test.go:315: TimescaleDB version: PostgreSQL 17.5 on aarch64-unknown-linux-musl, compiled by gcc (Alpine 14.2.0) 14.2.0, 64-bit
=== RUN   TestDatabaseConnections/ClickHouse
    benchmark_test.go:329: ClickHouse version: 25.5.1.2782
--- PASS: TestDatabaseConnections (0.05s)
    --- PASS: TestDatabaseConnections/TimescaleDB (0.04s)
    --- PASS: TestDatabaseConnections/ClickHouse (0.01s)
goos: darwin
goarch: arm64
pkg: db-bench
cpu: Apple M4
BenchmarkTimescaleDB_SingleInsert
BenchmarkTimescaleDB_SingleInsert-10              106872            108455 ns/op             283 B/op          8 allocs/op
BenchmarkTimescaleDB_BatchInsert
BenchmarkTimescaleDB_BatchInsert-10                  121         109814288 ns/op          261862 B/op       8718 allocs/op
BenchmarkClickHouse_SingleInsert
BenchmarkClickHouse_SingleInsert-10                21093            766671 ns/op           29282 B/op        257 allocs/op
BenchmarkClickHouse_AsyncInsert
BenchmarkClickHouse_AsyncInsert-10                110973            101958 ns/op            2276 B/op         45 allocs/op
BenchmarkClickHouse_BatchInsert
BenchmarkClickHouse_BatchInsert-10                  9183           1974473 ns/op          368096 B/op       6161 allocs/op
BenchmarkClickHouse_AsyncBatchInsert
BenchmarkClickHouse_AsyncBatchInsert-10            12108           1019599 ns/op          363449 B/op       6104 allocs/op
PASS
ok      db-bench        113.354s
```

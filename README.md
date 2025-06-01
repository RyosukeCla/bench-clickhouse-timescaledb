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
$ make benchmark                                                                                                                      3:30:08
go test -v -bench=. -benchmem -benchtime=10s
=== RUN   TestDatabaseConnections
=== RUN   TestDatabaseConnections/TimescaleDB
    benchmark_test.go:315: TimescaleDB version: PostgreSQL 17.5 on aarch64-unknown-linux-musl, compiled by gcc (Alpine 14.2.0) 14.2.0, 64-bit
=== RUN   TestDatabaseConnections/ClickHouse
    benchmark_test.go:329: ClickHouse version: 25.5.1.2782
--- PASS: TestDatabaseConnections (0.19s)
    --- PASS: TestDatabaseConnections/TimescaleDB (0.16s)
    --- PASS: TestDatabaseConnections/ClickHouse (0.03s)
goos: darwin
goarch: arm64
pkg: db-bench
cpu: Apple M4
BenchmarkTimescaleDB_SingleInsert
BenchmarkTimescaleDB_SingleInsert-10               40560            274616 ns/op             285 B/op          8 allocs/op
BenchmarkTimescaleDB_BatchInsert
BenchmarkTimescaleDB_BatchInsert-10                   46         219877293 ns/op          262047 B/op       8728 allocs/op
BenchmarkClickHouse_SingleInsert
BenchmarkClickHouse_SingleInsert-10                10000           1715473 ns/op           29585 B/op        260 allocs/op
BenchmarkClickHouse_AsyncInsert
BenchmarkClickHouse_AsyncInsert-10                   730          18791043 ns/op           29744 B/op        312 allocs/op
BenchmarkClickHouse_BatchInsert
BenchmarkClickHouse_BatchInsert-10                  3264           5848665 ns/op          363626 B/op       6174 allocs/op
BenchmarkClickHouse_AsyncBatchInsert
BenchmarkClickHouse_AsyncBatchInsert-10              182          63022697 ns/op          358177 B/op       6112 allocs/op
PASS
ok      db-bench        95.781s
```
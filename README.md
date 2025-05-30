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
$ make clean
```

## Result

```sh
$ make benchmark                                                                                  22:38:21
go test -v -bench=. -benchmem -benchtime=10s
=== RUN   TestDatabaseConnections
=== RUN   TestDatabaseConnections/TimescaleDB
    benchmark_test.go:315: TimescaleDB version: PostgreSQL 17.5 on aarch64-unknown-linux-musl, compiled by gcc (Alpine 14.2.0) 14.2.0, 64-bit
=== RUN   TestDatabaseConnections/ClickHouse
    benchmark_test.go:329: ClickHouse version: 25.5.1.2782
--- PASS: TestDatabaseConnections (0.04s)
    --- PASS: TestDatabaseConnections/TimescaleDB (0.03s)
    --- PASS: TestDatabaseConnections/ClickHouse (0.01s)
goos: darwin
goarch: arm64
pkg: db-bench
cpu: Apple M4
BenchmarkTimescaleDB_SingleInsert
BenchmarkTimescaleDB_SingleInsert-10              113576            117738 ns/op             282 B/op          8 allocs/op
BenchmarkTimescaleDB_BatchInsert
BenchmarkTimescaleDB_BatchInsert-10                  120         109456049 ns/op          261903 B/op       8721 allocs/op
BenchmarkClickHouse_SingleInsert
BenchmarkClickHouse_SingleInsert-10                26490            595812 ns/op           29147 B/op        256 allocs/op
BenchmarkClickHouse_AsyncInsert
BenchmarkClickHouse_AsyncInsert-10                133556             97104 ns/op            2268 B/op         45 allocs/op
BenchmarkClickHouse_BatchInsert
BenchmarkClickHouse_BatchInsert-10                  7072           1884508 ns/op          371176 B/op       6165 allocs/op
BenchmarkClickHouse_AsyncBatchInsert
BenchmarkClickHouse_AsyncBatchInsert-10            13873            843817 ns/op          363832 B/op       6107 allocs/op
PASS
ok      db-bench        115.421s

```

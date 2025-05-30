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
$ make benchmark                                                                                                                                 2:47:59
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
BenchmarkTimescaleDB_SingleInsert-10              110136            118593 ns/op             282 B/op          8 allocs/op
BenchmarkTimescaleDB_BatchInsert
BenchmarkTimescaleDB_BatchInsert-10                  118         102576179 ns/op          261918 B/op       8725 allocs/op
BenchmarkClickHouse_SingleInsert
BenchmarkClickHouse_SingleInsert-10                26350            630428 ns/op           29221 B/op        256 allocs/op
BenchmarkClickHouse_AsyncInsert
BenchmarkClickHouse_AsyncInsert-10                146108            101646 ns/op            2270 B/op         45 allocs/op
BenchmarkClickHouse_BatchInsert
BenchmarkClickHouse_BatchInsert-10                  5634           1874965 ns/op          366506 B/op       6168 allocs/op
BenchmarkClickHouse_AsyncBatchInsert
BenchmarkClickHouse_AsyncBatchInsert-10            14215            816058 ns/op          363827 B/op       6111 allocs/op
PASS
ok      db-bench        113.281s
```

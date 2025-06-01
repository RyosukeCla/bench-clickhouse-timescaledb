package main

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// TestData represents the data structure for our benchmark
type TestData struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	UserID    int64     `json:"user_id"`
	Value     float64   `json:"value"`
	Status    string    `json:"status"`
}

var (
	// Database connections
	timescaleConn  *pgxpool.Pool
	clickhouseConn driver.Conn

	// Test data
	testRecords []TestData
	numRecords  = 10000
)

func init() {
	// Generate test data
	testRecords = generateTestData(numRecords)
}

func generateTestData(count int) []TestData {
	rand.Seed(time.Now().UnixNano())
	data := make([]TestData, count)
	statuses := []string{"active", "inactive", "pending", "completed"}

	for i := 0; i < count; i++ {
		data[i] = TestData{
			ID:        int64(i + 1),
			Timestamp: time.Now().Add(-time.Duration(rand.Intn(86400)) * time.Second),
			UserID:    rand.Int63n(10000),
			Value:     rand.Float64() * 1000,
			Status:    statuses[rand.Intn(len(statuses))],
		}
	}
	return data
}

func setupTimescaleDB() error {
	config, err := pgxpool.ParseConfig("postgres://postgres:password@localhost:5432/benchmark?sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = 20
	config.MinConns = 5

	timescaleConn, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to connect to TimescaleDB: %w", err)
	}

	// Create hypertable
	_, err = timescaleConn.Exec(context.Background(), `
		DROP TABLE IF EXISTS test_data;
		CREATE TABLE test_data (
			id BIGINT NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL,
			user_id BIGINT NOT NULL,
			value DOUBLE PRECISION NOT NULL,
			status TEXT NOT NULL
		);
		SELECT create_hypertable('test_data', 'timestamp');
		CREATE INDEX ON test_data (user_id, timestamp DESC);
	`)

	return err
}

func setupClickHouse() error {
	var err error
	clickhouseConn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "benchmark",
			Username: "default",
			Password: "password",
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 30 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	// Drop table if exists
	err = clickhouseConn.Exec(context.Background(), `DROP TABLE IF EXISTS test_data`)
	if err != nil {
		return fmt.Errorf("failed to drop table: %w", err)
	}

	// Create table
	err = clickhouseConn.Exec(context.Background(), `
		CREATE TABLE test_data (
			id Int64,
			timestamp DateTime64(3),
			user_id Int64,
			value Float64,
			status String
		) ENGINE = MergeTree()
		ORDER BY (user_id, timestamp)
		PARTITION BY toYYYYMM(timestamp)
		SETTINGS index_granularity = 8192`)

	return err
}

func BenchmarkTimescaleDB_SingleInsert(b *testing.B) {
	if err := setupTimescaleDB(); err != nil {
		b.Fatalf("Failed to setup TimescaleDB: %v", err)
	}
	defer timescaleConn.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			record := testRecords[i%len(testRecords)]
			_, err := timescaleConn.Exec(context.Background(),
				"INSERT INTO test_data (id, timestamp, user_id, value, status) VALUES ($1, $2, $3, $4, $5)",
				record.ID, record.Timestamp, record.UserID, record.Value, record.Status)
			if err != nil {
				b.Errorf("Insert failed: %v", err)
			}
			i++
		}
	})
}

func BenchmarkTimescaleDB_BatchInsert(b *testing.B) {
	if err := setupTimescaleDB(); err != nil {
		b.Fatalf("Failed to setup TimescaleDB: %v", err)
	}
	defer timescaleConn.Close()

	batchSize := 1000
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tx, err := timescaleConn.Begin(context.Background())
		if err != nil {
			b.Fatalf("Failed to begin transaction: %v", err)
		}

		for j := 0; j < batchSize; j++ {
			record := testRecords[j%len(testRecords)]
			_, err := tx.Exec(context.Background(),
				"INSERT INTO test_data (id, timestamp, user_id, value, status) VALUES ($1, $2, $3, $4, $5)",
				record.ID, record.Timestamp, record.UserID, record.Value, record.Status)
			if err != nil {
				tx.Rollback(context.Background())
				b.Errorf("Insert failed: %v", err)
				break
			}
		}

		if err := tx.Commit(context.Background()); err != nil {
			b.Errorf("Commit failed: %v", err)
		}
	}
}

func BenchmarkClickHouse_SingleInsert(b *testing.B) {
	if err := setupClickHouse(); err != nil {
		b.Fatalf("Failed to setup ClickHouse: %v", err)
	}
	defer clickhouseConn.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			record := testRecords[i%len(testRecords)]
			err := clickhouseConn.Exec(context.Background(),
				"INSERT INTO test_data (id, timestamp, user_id, value, status) VALUES (?, ?, ?, ?, ?)",
				record.ID, record.Timestamp, record.UserID, record.Value, record.Status)
			if err != nil {
				b.Errorf("Insert failed: %v", err)
			}
			i++
		}
	})
}

func BenchmarkClickHouse_AsyncInsert(b *testing.B) {
	if err := setupClickHouse(); err != nil {
		b.Fatalf("Failed to setup ClickHouse: %v", err)
	}
	defer clickhouseConn.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			record := testRecords[i%len(testRecords)]
			// Use context with settings for async insert
			ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
				"async_insert":          1,
				"wait_for_async_insert": 1,
			}))
			err := clickhouseConn.Exec(ctx,
				"INSERT INTO test_data (id, timestamp, user_id, value, status) VALUES (?, ?, ?, ?, ?)",
				record.ID, record.Timestamp, record.UserID, record.Value, record.Status)
			if err != nil {
				b.Errorf("Async insert failed: %v", err)
			}
			i++
		}
	})
}

func BenchmarkClickHouse_BatchInsert(b *testing.B) {
	if err := setupClickHouse(); err != nil {
		b.Fatalf("Failed to setup ClickHouse: %v", err)
	}
	defer clickhouseConn.Close()

	batchSize := 1000
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		batch, err := clickhouseConn.PrepareBatch(context.Background(),
			"INSERT INTO test_data (id, timestamp, user_id, value, status)")
		if err != nil {
			b.Fatalf("Failed to prepare batch: %v", err)
		}

		for j := 0; j < batchSize; j++ {
			record := testRecords[j%len(testRecords)]
			err := batch.Append(record.ID, record.Timestamp, record.UserID, record.Value, record.Status)
			if err != nil {
				b.Errorf("Failed to append to batch: %v", err)
				break
			}
		}

		if err := batch.Send(); err != nil {
			b.Errorf("Failed to send batch: %v", err)
		}
	}
}

func BenchmarkClickHouse_AsyncBatchInsert(b *testing.B) {
	if err := setupClickHouse(); err != nil {
		b.Fatalf("Failed to setup ClickHouse: %v", err)
	}
	defer clickhouseConn.Close()

	batchSize := 1000
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Use context with settings for async insert
		ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
			"async_insert":          1,
			"wait_for_async_insert": 1,
		}))
		batch, err := clickhouseConn.PrepareBatch(ctx,
			"INSERT INTO test_data (id, timestamp, user_id, value, status)")
		if err != nil {
			b.Fatalf("Failed to prepare async batch: %v", err)
		}

		for j := 0; j < batchSize; j++ {
			record := testRecords[j%len(testRecords)]
			err := batch.Append(record.ID, record.Timestamp, record.UserID, record.Value, record.Status)
			if err != nil {
				b.Errorf("Failed to append to async batch: %v", err)
				break
			}
		}

		if err := batch.Send(); err != nil {
			b.Errorf("Failed to send async batch: %v", err)
		}
	}
}

// Utility function to wait for databases to be ready
func TestDatabaseConnections(t *testing.T) {
	t.Run("TimescaleDB", func(t *testing.T) {
		if err := setupTimescaleDB(); err != nil {
			t.Fatalf("Failed to connect to TimescaleDB: %v", err)
		}
		defer timescaleConn.Close()

		var version string
		err := timescaleConn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
		if err != nil {
			t.Fatalf("Failed to query TimescaleDB: %v", err)
		}
		t.Logf("TimescaleDB version: %s", version)
	})

	t.Run("ClickHouse", func(t *testing.T) {
		if err := setupClickHouse(); err != nil {
			t.Fatalf("Failed to connect to ClickHouse: %v", err)
		}
		defer clickhouseConn.Close()

		var version string
		err := clickhouseConn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
		if err != nil {
			t.Fatalf("Failed to query ClickHouse: %v", err)
		}
		t.Logf("ClickHouse version: %s", version)
	})
}

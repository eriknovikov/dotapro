package main

import (
	"time"

	"github.com/rs/zerolog/log"
)

// Metrics tracks scraping performance statistics
type Metrics struct {
	// Overall metrics
	StartTime          time.Time
	EndTime            time.Time
	TotalMatches       int
	TotalBatches       int
	SuccessfulBatches  int
	FailedBatches      int
	SkippedBatches     int

	// HTTP request metrics
	HTTPRequests       int
	HTTPSuccesses      int
	HTTPFailures       int
	HTTPRetryAttempts  int

	// Database metrics
	DBInserts          int
	DBInsertErrors     int

	// Timing metrics
	TotalFetchTime     time.Duration
	TotalProcessTime   time.Duration
	TotalDBTime        time.Duration
}

// NewMetrics creates a new Metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		StartTime: time.Now(),
	}
}

// StartBatch begins tracking a new batch
func (m *Metrics) StartBatch() {
	m.TotalBatches++
}

// RecordBatchSuccess records a successfully processed batch
func (m *Metrics) RecordBatchSuccess(matchCount int) {
	m.SuccessfulBatches++
	m.TotalMatches += matchCount
}

// RecordBatchFailure records a failed batch
func (m *Metrics) RecordBatchFailure() {
	m.FailedBatches++
}

// RecordBatchSkip records a skipped batch
func (m *Metrics) RecordBatchSkip() {
	m.SkippedBatches++
}

// RecordHTTPRequest records an HTTP request
func (m *Metrics) RecordHTTPRequest() {
	m.HTTPRequests++
}

// RecordHTTPSuccess records a successful HTTP request
func (m *Metrics) RecordHTTPSuccess() {
	m.HTTPSuccesses++
}

// RecordHTTPFailure records a failed HTTP request
func (m *Metrics) RecordHTTPFailure() {
	m.HTTPFailures++
}

// RecordHTTPRetry records a retry attempt
func (m *Metrics) RecordHTTPRetry() {
	m.HTTPRetryAttempts++
}

// RecordDBInsert records a database insert operation
func (m *Metrics) RecordDBInsert() {
	m.DBInserts++
}

// RecordDBInsertError records a database insert error
func (m *Metrics) RecordDBInsertError() {
	m.DBInsertErrors++
}

// RecordFetchTime records time spent fetching data
func (m *Metrics) RecordFetchTime(d time.Duration) {
	m.TotalFetchTime += d
}

// RecordProcessTime records time spent processing data
func (m *Metrics) RecordProcessTime(d time.Duration) {
	m.TotalProcessTime += d
}

// RecordDBTime records time spent on database operations
func (m *Metrics) RecordDBTime(d time.Duration) {
	m.TotalDBTime += d
}

// Finalize marks the end of scraping and calculates final metrics
func (m *Metrics) Finalize() {
	m.EndTime = time.Now()
}

// Log outputs the metrics to the logger
func (m *Metrics) Log() {
	m.Finalize()
	duration := m.EndTime.Sub(m.StartTime)

	log.Info().
		Dur("total_duration", duration).
		Int("total_matches", m.TotalMatches).
		Int("total_batches", m.TotalBatches).
		Int("successful_batches", m.SuccessfulBatches).
		Int("failed_batches", m.FailedBatches).
		Int("skipped_batches", m.SkippedBatches).
		Int("http_requests", m.HTTPRequests).
		Int("http_successes", m.HTTPSuccesses).
		Int("http_failures", m.HTTPFailures).
		Int("http_retries", m.HTTPRetryAttempts).
		Int("db_inserts", m.DBInserts).
		Int("db_insert_errors", m.DBInsertErrors).
		Dur("total_fetch_time", m.TotalFetchTime).
		Dur("total_process_time", m.TotalProcessTime).
		Dur("total_db_time", m.TotalDBTime).
		Float64("matches_per_second", float64(m.TotalMatches)/duration.Seconds()).
		Msg("scraping metrics")
}

// BatchTimer tracks timing for a single batch operation
type BatchTimer struct {
	metrics    *Metrics
	startTime  time.Time
	operation  string
}

// NewBatchTimer creates a new batch timer
func NewBatchTimer(metrics *Metrics, operation string) *BatchTimer {
	return &BatchTimer{
		metrics:   metrics,
		startTime: time.Now(),
		operation: operation,
	}
}

// Stop stops the timer and records the duration
func (bt *BatchTimer) Stop() {
	duration := time.Since(bt.startTime)

	switch bt.operation {
	case "fetch":
		bt.metrics.RecordFetchTime(duration)
	case "process":
		bt.metrics.RecordProcessTime(duration)
	case "db":
		bt.metrics.RecordDBTime(duration)
	}

	log.Debug().
		Str("operation", bt.operation).
		Dur("duration", duration).
		Msg("batch operation completed")
}

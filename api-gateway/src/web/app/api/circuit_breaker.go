package api

import (
	"log"
	"time"

	"github.com/sony/gobreaker/v2"
)

// CircuitBreakerConfig vsebuje konfiguracijo za circuit breaker
type CircuitBreakerConfig struct {
	MaxRequests     uint32
	Interval        time.Duration
	Timeout         time.Duration
	MinRequestCount uint32
	FailureRatio    float64
}

// DefaultConfig vrne privzeto konfiguracijo circuit breaker-ja
func DefaultConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxRequests:     3,
		Interval:        10 * time.Second,
		Timeout:         30 * time.Second,
		MinRequestCount: 10,
		FailureRatio:    0.6,
	}
}

// NewCircuitBreaker ustvari nov circuit breaker s podano konfiguracijo
func NewCircuitBreaker[T any](name string, config CircuitBreakerConfig) *gobreaker.CircuitBreaker[T] {
	return gobreaker.NewCircuitBreaker[T](gobreaker.Settings{
		Name:        name,
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= config.MinRequestCount && failureRatio >= config.FailureRatio
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit Breaker %s state changed from %s to %s", name, from, to)
		},
	})
}

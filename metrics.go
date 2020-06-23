package metrics

import (
	"github.com/go-redis/redis/v7"
)

type Metrics struct {
	client redis.UniversalClient

	System   *SystemMetrics
	Internal *InternalMetrics
	Time     *AvgTimeMetric
}

func NewMetrics(client redis.UniversalClient, interval int, globalPrefix string) *Metrics {
	if globalPrefix == "" {
		globalPrefix = "metrics:"
	}

	systemMetrics := NewSystemMetrics()
	internalMetrics := NewInternalMetrics(client, interval, globalPrefix)
	timeMetrics := NewAvgTimeMetric(client, interval, globalPrefix)

	return &Metrics{
		client:   client,
		System:   systemMetrics,
		Internal: internalMetrics,
		Time:     timeMetrics,
	}
}

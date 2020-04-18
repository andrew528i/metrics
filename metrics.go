package metrics

import (
	"github.com/go-redis/redis/v7"
)

type Metrics struct {
	client *redis.Client

	System   *SystemMetrics
	Internal *InternalMetrics
	Time     *AvgTimeMetric
}

func NewMetrics(client *redis.Client, interval int) *Metrics {
	systemMetrics := NewSystemMetrics()
	internalMetrics := NewInternalMetrics(client, interval)
	timeMetrics := NewAvgTimeMetric(client, interval)

	return &Metrics{
		client:   client,
		System:   systemMetrics,
		Internal: internalMetrics,
		Time:     timeMetrics,
	}
}

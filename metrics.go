package metrics

import (
	"github.com/go-redis/redis/v7"
)

type Metrics struct {
	client *redis.ClusterClient

	System   *SystemMetrics
	Internal *InternalMetrics
	Time     *AvgTimeMetric
}

func NewMetrics(client *redis.ClusterClient) *Metrics {
	systemMetrics := NewSystemMetrics()
	internalMetrics := NewInternalMetrics(client)
	timeMetrics := NewAvgTimeMetric(client, 3)

	return &Metrics{
		client:   client,
		System:   systemMetrics,
		Internal: internalMetrics,
		Time:     timeMetrics,
	}
}

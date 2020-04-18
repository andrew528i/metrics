package metrics

import "github.com/go-redis/redis/v7"

type InternalMetrics struct {
	Counter    *CounterMetric
	AvgCounter *AvgCounterMetric
}

func NewInternalMetrics(client *redis.ClusterClient) *InternalMetrics {
	counter := NewCounterMetric(client, 3)
	avgCounter := NewAvgCounterMetric(client, 3)

	return &InternalMetrics{Counter: counter, AvgCounter: avgCounter}
}

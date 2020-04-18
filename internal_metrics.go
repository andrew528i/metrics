package metrics

import "github.com/go-redis/redis/v7"

type InternalMetrics struct {
	Counter    *CounterMetric
	AvgCounter *AvgCounterMetric
}

func NewInternalMetrics(client *redis.Client, interval int) *InternalMetrics {
	counter := NewCounterMetric(client, interval)
	avgCounter := NewAvgCounterMetric(client, interval)

	return &InternalMetrics{Counter: counter, AvgCounter: avgCounter}
}

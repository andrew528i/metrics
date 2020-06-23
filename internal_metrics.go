package metrics

import "github.com/go-redis/redis/v7"

type InternalMetrics struct {
	Counter    *CounterMetric
	AvgCounter *AvgCounterMetric
}

func NewInternalMetrics(client redis.UniversalClient, interval int, globalPrefix string) *InternalMetrics {
	counter := NewCounterMetric(client, interval, globalPrefix)
	avgCounter := NewAvgCounterMetric(client, interval, globalPrefix)

	return &InternalMetrics{Counter: counter, AvgCounter: avgCounter}
}

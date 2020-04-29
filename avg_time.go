package metrics

import (
	"time"

	"github.com/go-redis/redis/v7"
)

type (
	AvgTimeMetric struct {
		interval   int // in minutes
		avgCounter *AvgCounterMetric
	}

	AvgTimeMetricTimer struct {
		names      []string
		avgCounter *AvgCounterMetric
		startedAt  int64
	}
)

func NewAvgTimeMetric(client redis.UniversalClient, interval int) *AvgTimeMetric {
	avgCounter := NewAvgCounterMetric(client, interval)

	return &AvgTimeMetric{avgCounter: avgCounter, interval: interval}
}

func (s AvgTimeMetric) Start(names ...string) *AvgTimeMetricTimer {
	timer := &AvgTimeMetricTimer{names: names, avgCounter: s.avgCounter}
	timer.Start()

	return timer
}

func (s *AvgTimeMetricTimer) Start() {
	s.startedAt = time.Now().UnixNano()
}

func (s AvgTimeMetricTimer) Stop() {
	stoppedAt := time.Now().UnixNano()
	dt := stoppedAt - s.startedAt

	for _, name := range s.names {
		// NOTE: second -> millisecond -> microsecond -> nanosecond
		s.avgCounter.Add(name, dt / 1000000)
	}
}

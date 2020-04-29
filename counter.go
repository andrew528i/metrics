package metrics

import (
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

type CounterMetric struct {
	client redis.UniversalClient

	interval int // in minutes
}

func NewCounterMetric(client redis.UniversalClient, interval int) *CounterMetric {
	return &CounterMetric{client: client, interval: interval}
}

func (s CounterMetric) Get(name string) (int, error) {
	now := time.Now()
	total := 0

	for i := 0; i < s.interval; i += 1 {
		key := getKeyByDate(now.Add(-1 * time.Minute * time.Duration(i)), "counter", name)
		exists, err := s.client.Exists(key).Result()
		if err != nil {
			return 0, err
		}

		if exists == 0 {
			continue
		}

		res, err := s.client.Get(key).Result()
		if err != nil {
			return 0, err
		}

		parsed, err := strconv.Atoi(res)
		if err != nil {
			return 0, nil
		}

		total += parsed
	}

	return total, nil
}

func (s CounterMetric) Increase(name string) error {
	key := getKey("counter", name)
	result, err := s.client.Exists(key).Result()
	if err != nil {
		return err
	}

	if result == 0 {
		s.client.SetNX(key, 1, time.Minute * time.Duration(s.interval) * 2)
	} else {
		s.client.Incr(key)
	}

	return nil
}


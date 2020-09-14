package metrics

import (
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

type CounterMetric struct {
	client redis.UniversalClient

	interval     int // in minutes
	globalPrefix string
}

func NewCounterMetric(client redis.UniversalClient, interval int, globalPrefix string) *CounterMetric {
	return &CounterMetric{
		client:       client,
		interval:     interval,
		globalPrefix: globalPrefix,
	}
}

func (s CounterMetric) Get(name string) (int, error) {
	now := time.Now()
	total := 0

	for i := 0; i < s.interval; i += 1 {
		key := getKeyByDate(now.Add(-1*time.Minute*time.Duration(i)), s.globalPrefix+"counter", name)
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

func (s CounterMetric) Increase(name string, vals ...int64) error {
	key := getKey(s.globalPrefix+"counter", name)
	result, err := s.client.Exists(key).Result()
	if err != nil {
		return err
	}
	var val int64 = 1
	if len(vals) > 0 {
		val = 0
		for _, v := range vals {
			val += v
		}
	}
	if result == 0 {
		s.client.SetNX(key, val, time.Minute*time.Duration(s.interval)*2)
	} else {
		s.client.IncrBy(key, val)
	}

	return nil
}

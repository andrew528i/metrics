package metrics

import (
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

type AvgCounterMetric struct {
	client *redis.ClusterClient

	interval int // in minutes
}

func NewAvgCounterMetric(client *redis.ClusterClient, interval int) *AvgCounterMetric {
	return &AvgCounterMetric{client: client, interval: interval}
}

func (s AvgCounterMetric) Add(name string, ms int64) error {
	key := getKey("avg_counter", name)
	exists, err := s.client.Exists(key).Result()
	if err != nil {
		return err
	}

	_, err = s.client.LPush(key, ms).Result()
	if err != nil {
		return err
	}

	if exists == 0 {
		s.client.Expire(key, time.Minute * time.Duration(s.interval) * 2)
	}

	return nil
}

func (s AvgCounterMetric) GetAvg(name string) (int64, error) {
	now := time.Now()
	values := make([]int64, 0)
	valueSum := int64(0)
	valueCount := int64(0)

	for i := 0; i < s.interval; i += 1 {
		key := getKeyByDate(
			now.Add(-1 * time.Minute * time.Duration(i)), "avg_counter", name)
		exists, err := s.client.Exists(key).Result()
		if err != nil {
			return 0, err
		}

		if exists == 0 {
			continue
		}

		res, err := s.client.LRange(key, 0, -1).Result()
		if err != nil {
			return 0, err
		}

		for _, value := range res {
			parsed, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return 0, nil
			}

			values = append(values, parsed)
			valueSum += parsed
			valueCount += 1
		}
	}

	if valueCount == 0 {
		return 0, nil
	}

	average := valueSum / valueCount

	return average, nil
}

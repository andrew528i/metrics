# Installation

```console
go get -u github.com/andrew528i/metrics@v1.0.1
```

# Usage

```go
package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/andrew528i/metrics"
	"github.com/go-redis/redis/v7"
)

// some heavy func, i.e. photo processing or smth else
func heavy() {
	dur := time.Millisecond * time.Duration(rand.Intn(1000))
	time.Sleep(dur)
}

func main() {
	// initialize redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "",
	})

	interval := 3 // in minutes
	m := metrics.NewMetrics(client, interval)

	// counter metric
	for i := 0; i < 10; i += 1 {
		m.Internal.Counter.Increase("my_value")
	}

	// average counter metric
	for i := 0; i < 10; i += 1 {
		v := int64(rand.Intn(100))
		m.Internal.AvgCounter.Add("my_avg_value", v)
	}

	// execution time metric
	var wg sync.WaitGroup
	for i := 0; i < 25; i += 1 {
		wg.Add(1)

		go func() {
			timer := m.Time.Start("heavy_avg_ts")
			heavy()
			timer.Stop()
			wg.Done()
		}()
	}

	wg.Wait()

	// run metrics http server
	srv := metrics.NewServer(m)
	srv.Run()
}
```

After running this code you can access metrics above with:
```text
http://localhost:6767/internal/my_value
http://localhost:6767/internal/my_avg_value
http://localhost:6767/internal/heavy_avg_ts
```

Also there are some system metrics:
```text
http://localhost:6767/system/cpu  // in percents
http://localhost:6767/system/ram  // in megabytes
http://localhost:6767/system/disk // disk usage in percents
```

There are 3 types of internal metrics for now:

## Counter

You can count new users, orders, etc in last N minutes.

```
m.Internal.Counter.Increase("new_users")
```

`http://localhost:6767/internal/new_users`

## Average Counter

Can be used for measuring average order sum in last N minutes

```
m.Internal.AvgCounter.Add("order_sum", order.sum)
```

`http://localhost:6767/internal/order_sum`

## Average execution time

An average HTTP/GRPC/NATS request execution can be measured:

```
timer := m.Time.Start("process_photo_ts")
process_photo()
timer.Stop()
```

`http://localhost:6767/internal/process_photo_ts`

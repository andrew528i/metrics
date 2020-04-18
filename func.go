package metrics

import (
	"fmt"
	"time"
)

func getKeyByDate(time time.Time, prefix, name string) string {
	formattedDate := time.Format("200601021504")

	// TODO: remove this hardcoded global prefix
	return fmt.Sprintf("metrics:%s:%s:%s", prefix, name, formattedDate)
}

func getKey(prefix, name string) string {
	now := time.Now()

	return getKeyByDate(now, prefix, name)
}


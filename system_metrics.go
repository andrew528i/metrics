package metrics

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
)

type SystemMetrics struct {}

func NewSystemMetrics() *SystemMetrics {
	return &SystemMetrics{}
}

func (SystemMetrics) CPU() float64 {
	stat, _ := cpu.Percent(time.Second, true)
	averageLoad := .0

	for _, load := range stat {
		averageLoad += load
	}

	return averageLoad / float64(len(stat))
}

func (SystemMetrics) RAM() float64 {
	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	return float64(m.Alloc) / 1024 / 1024
}

func (SystemMetrics) Disk() float64 {
	stat, _ := disk.Usage("/")
	used := stat.Total - stat.Free

	return 100 * float64(used) / float64(stat.Total)
}

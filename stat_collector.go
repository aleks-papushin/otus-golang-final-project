package main

import (
	"time"
)

type StatCollector interface {
	getStatSnapshot() (string, error)
	parseDate(string) time.Time
	parseLastSecLoadAverage(string) float32
	parseCpuUsage(string) CpuUsage
}

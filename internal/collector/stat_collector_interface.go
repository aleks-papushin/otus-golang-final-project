package collector

import (
	"time"

	"github.com/aleks-papushin/system-monitor/internal/models"
)

const (
	minInterval  = 5
	maxSnapshots = 24 * 60 * 60
)

type StatCollector interface {
	GetStatSnapshot() (string, error)
	ParseDate(string) time.Time
	ParseLastSecLoadAverage(string) float32
	ParseCpuUsage(string) models.CpuUsage
}

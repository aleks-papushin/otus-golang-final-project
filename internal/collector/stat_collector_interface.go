package collector

import (
	"time"

	"github.com/aleks-papushin/system-monitor/internal/models"
)

const (
	defaultCollectingInterval = 5 * time.Second
)

type StatCollector interface {
	GetStatSnapshot() (string, error)
	ParseDate(string) time.Time
	ParseLastSecLoadAverage(string) float32
	ParseCpuUsage(string) models.CpuUsage
}

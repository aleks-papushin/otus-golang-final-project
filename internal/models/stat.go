package models

import (
	"fmt"
	"time"
)

type Stat struct {
	LoadAverage float32
	CPUUsage    CPUUsage
	Time        time.Time
}

type CPUUsage struct {
	UserUsage float32
	SysUsage  float32
	Idle      float32
}

func (s Stat) String() string {
	return fmt.Sprintf("Time: %s, LoadAverage: %.2f, UserUsage: %.2f%%, SysUsage: %.2f%%, Idle: %.2f%%",
		s.Time.Format(time.RFC3339), s.LoadAverage, s.CPUUsage.UserUsage, s.CPUUsage.SysUsage, s.CPUUsage.Idle)
}

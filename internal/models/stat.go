package models

import (
	"fmt"
	"time"
)

type Stat struct {
	LoadAverage float32
	CpuUsage    CpuUsage
	Time        time.Time
}

type CpuUsage struct {
	UserUsage float32
	SysUsage  float32
	Idle      float32
}

func (s Stat) String() string {
	return fmt.Sprintf("Time: %s, LoadAverage: %.2f, UserUsage: %.2f%%, SysUsage: %.2f%%, Idle: %.2f%%",
		s.Time.Format(time.RFC3339), s.LoadAverage, s.CpuUsage.UserUsage, s.CpuUsage.SysUsage, s.CpuUsage.Idle)
}

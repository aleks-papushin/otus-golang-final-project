package models

import "time"

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

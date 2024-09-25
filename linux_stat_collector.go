package main

import (
	"strconv"
	"strings"
	"time"
)

type LinuxStatCollector struct {
	executor CommandExecutor
}

func (c *LinuxStatCollector) getStatSnapshot() (string, error) {
	outputBytes, err := c.executor.Execute("top", "-b", "-n 1")
	if err != nil {
		return "", err
	}
	return string(outputBytes), nil
}

func (c *LinuxStatCollector) parseDate(date string) time.Time {
	layout := "2006/01/02 15:04:05"
	parsedDate, _ := time.Parse(layout, date)
	return parsedDate
}

func (c *LinuxStatCollector) parseLastSecLoadAverage(line string) float32 {
	laString := strings.TrimSuffix(strings.Split(line, " ")[2], ",")
	la, _ := strconv.ParseFloat(laString, 32)
	return float32(la)
}

func (c *LinuxStatCollector) parseCpuUsage(line string) CpuUsage {
	cpuSlice := strings.Split(line, " ")
	userUsage := strings.TrimSuffix(cpuSlice[1], "%")
	sysUsage := strings.TrimSuffix(cpuSlice[3], "%")
	idle := strings.TrimSuffix(cpuSlice[7], "%")

	u, _ := strconv.ParseFloat(userUsage, 32)
	s, _ := strconv.ParseFloat(sysUsage, 32)
	i, _ := strconv.ParseFloat(idle, 32)

	return CpuUsage{
		userUsage: float32(u),
		sysUsage:  float32(s),
		idle:      float32(i),
	}
}

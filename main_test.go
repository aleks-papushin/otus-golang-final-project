package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseLastSecLoadAverage(t *testing.T) {
	collector := &MacOSStatCollector{}
	for _, tc := range []struct {
		Name      string
		TestData  string
		ExpResult float32
	}{
		{
			Name:      "Regular load average entry",
			TestData:  "Load Avg: 1.23, 0.98, 0.76",
			ExpResult: 1.23,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult := collector.parseLastSecLoadAverage(tc.TestData)
			require.Equal(t, tc.ExpResult, actualResult)
		})
	}
}

func TestParseCpuUsage(t *testing.T) {
	collector := &MacOSStatCollector{}
	for _, tc := range []struct {
		Name      string
		TestData  string
		ExpResult CpuUsage
	}{
		{
			Name:     "Regular CPU usage entry",
			TestData: "CPU usage: 2.36% user, 5.45% sys, 92.18% idle",
			ExpResult: CpuUsage{
				userUsage: 2.36,
				sysUsage:  5.45,
				idle:      92.18,
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult := collector.parseCpuUsage(tc.TestData)
			require.Equal(t, tc.ExpResult, actualResult)
		})
	}
}

func TestParseDate(t *testing.T) {
	collector := &MacOSStatCollector{}
	for _, tc := range []struct {
		Name      string
		TestData  string
		ExpResult time.Time
	}{
		{
			Name:      "Regular date entry",
			TestData:  "2023/10/05 14:30:00",
			ExpResult: time.Date(2023, 10, 5, 14, 30, 0, 0, time.UTC),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult := collector.parseDate(tc.TestData)
			require.Equal(t, tc.ExpResult, actualResult)
		})
	}
}

func TestGetStatSnapshot(t *testing.T) {
	mockOutput := []byte("bunch of stats from stat collecting utility")
	mockExecutor := &MockCommandExecutor{
		Output: mockOutput,
		Err:    nil,
	}
	collector := &MacOSStatCollector{
		executor: mockExecutor,
	}

	result, err := collector.getStatSnapshot()
	require.NoError(t, err)
	require.Equal(t, string(mockOutput), result)
}

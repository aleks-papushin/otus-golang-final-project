package collector

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aleks-papushin/system-monitor/internal/models"
)

type MacOSStatCollector struct {
	Executor  CommandExecutor
	snapshots []*models.Stat
}

var (
	instance *MacOSStatCollector
	once     sync.Once
)

func GetMacOSStatCollector() *MacOSStatCollector {
	once.Do(func() {
		instance = &MacOSStatCollector{
			Executor:  &RealCommandExecutor{},
			snapshots: make([]*models.Stat, 0),
		}
		go instance.startCollecting()
	})
	return instance
}

func (c *MacOSStatCollector) GetStatSnapshot() (string, error) {
	outputBytes, err := c.Executor.Execute("top", "-l 1", "-n 0")
	if err != nil {
		return "", err
	}
	return string(outputBytes), nil
}

func (c *MacOSStatCollector) ParseDate(date string) time.Time {
	layout := "2006/01/02 15:04:05"
	parsedDate, _ := time.Parse(layout, date)
	return parsedDate
}

func (c *MacOSStatCollector) ParseLastSecLoadAverage(line string) float32 {
	laString := strings.TrimSuffix(strings.Split(line, " ")[2], ",")
	la, _ := strconv.ParseFloat(laString, 32)
	return float32(la)
}

func (c *MacOSStatCollector) ParseCpuUsage(line string) models.CpuUsage {
	cpuSlice := strings.Split(line, " ")
	userUsage := strings.TrimSuffix(cpuSlice[2], "%")
	sysUsage := strings.TrimSuffix(cpuSlice[4], "%")
	idle := strings.TrimSuffix(cpuSlice[6], "%")

	u, _ := strconv.ParseFloat(userUsage, 32)
	s, _ := strconv.ParseFloat(sysUsage, 32)
	i, _ := strconv.ParseFloat(idle, 32)

	return models.CpuUsage{
		UserUsage: float32(u),
		SysUsage:  float32(s),
		Idle:      float32(i),
	}
}

func (c *MacOSStatCollector) CollectMacOSStat(outputInterval, collectingInterval int) <-chan *models.Stat {
	collectingIntervalSec := time.Duration(collectingInterval) * time.Second
	avgStatChan := make(chan *models.Stat)

	go func() {
		collectT := time.NewTimer(collectingIntervalSec)
		<-collectT.C // wait collectingInterval seconds before taking first snapshot...

		outputT := time.NewTicker(time.Duration(outputInterval) * time.Second)
		for {
			timeEdge := time.Now().Add(-collectingIntervalSec)
			avgStatSnapshot := c.makeAvgSnapshotAfter(timeEdge)
			avgStatChan <- avgStatSnapshot
			<-outputT.C // ...then make new snapshot every outputInterval seconds
		}
	}()

	return avgStatChan
}

func (c *MacOSStatCollector) makeAvgSnapshotAfter(t time.Time) *models.Stat {
	laSum := float32(0.0)
	uCpuSum := float32(0.0)
	sCpuSum := float32(0.0)
	iCpuSum := float32(0.0)
	snapShotsCount := float32(0.0)
	for i := len(c.snapshots) - 1; i >= 0; i-- {
		s := c.snapshots[i]
		if s.Time.After(t) {
			snapShotsCount++
			laSum += s.LoadAverage
			uCpuSum += s.CpuUsage.UserUsage
			sCpuSum += s.CpuUsage.SysUsage
			iCpuSum += s.CpuUsage.Idle
		} else {
			break
		}
	}

	avgStatSnapshot := models.Stat{
		LoadAverage: laSum / snapShotsCount,
		CpuUsage: models.CpuUsage{
			UserUsage: uCpuSum / snapShotsCount,
			SysUsage:  sCpuSum / snapShotsCount,
			Idle:      iCpuSum / snapShotsCount,
		},
		Time: time.Now(),
	}
	return &avgStatSnapshot
}

func (c *MacOSStatCollector) startCollecting() {
	ticker := time.NewTicker(defaultCollectingInterval)

	go func() {
		for {
			select {
			case <-ticker.C:
				statSnapshot, err := c.GetStatSnapshot()
				statSlice := strings.Split(statSnapshot, "\n")
				if err != nil {
					errOut := fmt.Errorf("error occured on getting stat snapshot %w", err)
					fmt.Println(errOut)
				}
				t := c.ParseDate(statSlice[1])
				la := c.ParseLastSecLoadAverage(statSlice[2])
				cpu := c.ParseCpuUsage(statSlice[3])
				s := models.Stat{
					LoadAverage: la,
					CpuUsage:    cpu,
					Time:        t,
				}
				c.snapshots = append(c.snapshots, &s)
			}
		}
	}()
}

package collector

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aleks-papushin/system-monitor/config"
	"github.com/aleks-papushin/system-monitor/internal/models"
	ospackage "github.com/aleks-papushin/system-monitor/internal/os"
	"github.com/spf13/viper"
)

const (
	defaultCollectingInterval = 5 * time.Second
)

type OSStatCollector struct {
	Executor           CommandExecutor
	snapshots          []*models.Stat
	collectLoadAverage bool
	collectCPUUsage    bool
}

var (
	instance *OSStatCollector
	once     sync.Once
)

func GetMacOSStatCollector() *OSStatCollector {
	once.Do(func() {
		config.InitConfig()
		instance = &OSStatCollector{
			Executor:           &RealCommandExecutor{},
			snapshots:          make([]*models.Stat, 0),
			collectLoadAverage: viper.GetBool("collect.load_average"),
			collectCPUUsage:    viper.GetBool("collect.cpu_usage"),
		}
		go instance.startCollecting()
	})
	return instance
}

func (c *OSStatCollector) GetStatSnapshot() (string, error) {
	statCmd := ospackage.SysMonitorCmd
	outputBytes, err := c.Executor.Execute(statCmd[0], statCmd[1], statCmd[2])
	if err != nil {
		return "", err
	}
	return string(outputBytes), nil
}

func (c *OSStatCollector) ParseDate(date string) time.Time {
	layout := "2006/01/02 15:04:05"
	parsedDate, _ := time.Parse(layout, date)
	return parsedDate
}

func (c *OSStatCollector) ParseLastSecLoadAverage(line string) float32 {
	laString := strings.TrimSuffix(strings.Split(line, " ")[2], ",")
	la, _ := strconv.ParseFloat(laString, 32)
	return float32(la)
}

func (c *OSStatCollector) ParseCPUUsage(line string) models.CPUUsage {
	cpuSlice := strings.Split(line, " ")
	userUsage := strings.TrimSuffix(cpuSlice[2], "%")
	sysUsage := strings.TrimSuffix(cpuSlice[4], "%")
	idle := strings.TrimSuffix(cpuSlice[6], "%")

	u, _ := strconv.ParseFloat(userUsage, 32)
	s, _ := strconv.ParseFloat(sysUsage, 32)
	i, _ := strconv.ParseFloat(idle, 32)

	return models.CPUUsage{
		UserUsage: float32(u),
		SysUsage:  float32(s),
		Idle:      float32(i),
	}
}

func (c *OSStatCollector) CollectStat(outputInterval, collectingInterval int) <-chan *models.Stat {
	statInterval := time.Duration(collectingInterval) * time.Second
	avgStatChan := make(chan *models.Stat)

	go func() {
		collectT := time.NewTimer(statInterval)
		<-collectT.C // wait collectingInterval seconds before taking first snapshot...

		outputT := time.NewTicker(time.Duration(outputInterval) * time.Second)
		for {
			timeEdge := time.Now().Add(-statInterval)
			avgStatSnapshot := c.makeAvgSnapshotAfter(timeEdge)
			avgStatChan <- avgStatSnapshot
			<-outputT.C // ...then make new snapshot every outputInterval seconds
		}
	}()

	return avgStatChan
}

func (c *OSStatCollector) makeAvgSnapshotAfter(t time.Time) *models.Stat {
	laSum := float32(0.0)
	uCPUSum := float32(0.0)
	sCPUSum := float32(0.0)
	iCPUSum := float32(0.0)
	snapShotsCount := float32(0.0)
	for i := len(c.snapshots) - 1; i >= 0; i-- {
		s := c.snapshots[i]
		if s.Time.After(t) {
			snapShotsCount++
			laSum += s.LoadAverage
			uCPUSum += s.CPUUsage.UserUsage
			sCPUSum += s.CPUUsage.SysUsage
			iCPUSum += s.CPUUsage.Idle
		} else {
			break
		}
	}

	avgStatSnapshot := models.Stat{
		LoadAverage: laSum / snapShotsCount,
		CPUUsage: models.CPUUsage{
			UserUsage: uCPUSum / snapShotsCount,
			SysUsage:  sCPUSum / snapShotsCount,
			Idle:      iCPUSum / snapShotsCount,
		},
		Time: time.Now(),
	}
	return &avgStatSnapshot
}

func (c *OSStatCollector) startCollecting() {
	ticker := time.NewTicker(defaultCollectingInterval)

	go func() {
		for {
			<-ticker.C
			statSnapshot, err := c.GetStatSnapshot()
			statSlice := strings.Split(statSnapshot, "\n")
			if err != nil {
				errOut := fmt.Errorf("error occurred on getting stat snapshot %w", err)
				fmt.Println(errOut)
			}
			t := c.ParseDate(statSlice[1])
			var la float32
			var cpu models.CPUUsage

			if c.collectLoadAverage {
				la = c.ParseLastSecLoadAverage(statSlice[2])
			}
			if c.collectCPUUsage {
				cpu = c.ParseCPUUsage(statSlice[3])
			}
			s := models.Stat{
				LoadAverage: la,
				CPUUsage:    cpu,
				Time:        t,
			}
			c.snapshots = append(c.snapshots, &s)
		}
	}()
}

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
	Executor CommandExecutor
}

var (
	instance *MacOSStatCollector
	once     sync.Once
)

func GetMacOSStatCollector() *MacOSStatCollector {
	once.Do(func() {
		instance = &MacOSStatCollector{
			Executor: &RealCommandExecutor{},
		}
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

func (c *MacOSStatCollector) CollectMacOSStat(statChan chan *models.Stat, n, m int) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	wg.Add(1)
	c.startStatCollecting(n, m, statChan, &wg)

	wg.Wait()
}

func (c *MacOSStatCollector) startStatCollecting(n, m int, statChan chan *models.Stat, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				statSnapshot, err := c.GetStatSnapshot()
				statSlice := strings.Split(statSnapshot, "\n")
				if err != nil {
					errOut := fmt.Errorf("error occured attempting get load average %w", err)
					fmt.Println(errOut)
				}
				t := c.ParseDate(statSlice[1])
				la := c.ParseLastSecLoadAverage(statSlice[2])
				cpu := c.ParseCpuUsage(statSlice[3])

				stat := models.Stat{
					LoadAverage: la,
					CpuUsage:    cpu,
					Time:        t,
				}

				statChan <- &stat
			}
		}
	}()
}

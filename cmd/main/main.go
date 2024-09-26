package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aleks-papushin/system-monitor/internal/collector"
	"github.com/aleks-papushin/system-monitor/internal/models"
)

const (
	statCollectingInterval = 5
	maxM                   = 24 * 60 * 60
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	n := 5
	m := 15

	statChan := make(chan *models.Stat, m)
	c := &collector.MacOSStatCollector{}

	wg.Add(1)
	startStatCollecting(n, m, statChan, &wg, c)

	wg.Wait()
}

func startStatCollecting(n, m int, statChan chan *models.Stat, wg *sync.WaitGroup, c collector.StatCollector) {
	defer wg.Done()
	ticker := time.NewTicker(statCollectingInterval * time.Second)

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

	go func() {
		var stat *models.Stat
		for {
			select {
			case stat = <-statChan:
				fmt.Println("Received stat:", stat)
			}
		}
	}()
}

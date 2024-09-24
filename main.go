package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	statCollectingInterval = 5
	maxM                   = 24 * 60 * 60
)

type Stat struct {
	loadAverage float32
	cpuUsage    CpuUsage
	time        time.Time
}

type CpuUsage struct {
	userUsage float32
	sysUsage  float32
	idle      float32
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	n := 5
	m := 15

	statChan := make(chan *Stat, m)
	collector := &MacOSStatCollector{}

	wg.Add(1)
	startStatCollecting(n, m, statChan, &wg, collector)

	wg.Wait()
}

func startStatCollecting(n, m int, statChan chan *Stat, wg *sync.WaitGroup, collector StatCollector) {
	defer wg.Done()
	ticker := time.NewTicker(statCollectingInterval * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				statSnapshot, err := collector.getStatSnapshot()
				statSlice := strings.Split(statSnapshot, "\n")
				if err != nil {
					errOut := fmt.Errorf("error occured attempting get load average %w", err)
					fmt.Println(errOut)
				}
				t := collector.parseDate(statSlice[1])
				la := collector.parseLastSecLoadAverage(statSlice[2])
				cpu := collector.parseCpuUsage(statSlice[3])

				stat := Stat{
					loadAverage: la,
					cpuUsage:    cpu,
					time:        t,
				}

				statChan <- &stat
			}
		}
	}()

	go func() {
		var stat *Stat
		for {
			select {
			case stat = <-statChan:
				fmt.Println("Received stat:", stat)
			}
		}
	}()
}

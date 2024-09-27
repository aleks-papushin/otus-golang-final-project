package main

import (
	"sync"

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
	collector.StartStatCollecting(n, m, statChan, &wg, c)

	wg.Wait()
}

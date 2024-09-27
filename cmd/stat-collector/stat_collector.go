package main

import (
	"sync"

	"github.com/aleks-papushin/system-monitor/internal/collector"
	"github.com/aleks-papushin/system-monitor/internal/models"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	n := 5
	m := 15

	statChan := make(chan *models.Stat, m)
	c := &collector.MacOSStatCollector{}

	wg.Add(1)

	c.CollectMacOSStat(statChan, n, m)

	wg.Wait()
}

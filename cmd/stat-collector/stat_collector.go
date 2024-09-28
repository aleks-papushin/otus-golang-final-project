package main

import (
	"fmt"
	"sync"

	"github.com/aleks-papushin/system-monitor/internal/collector"
	"github.com/aleks-papushin/system-monitor/internal/models"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	n := 5
	m := 15

	c := &collector.MacOSStatCollector{}

	wg.Add(1)

	statChan := c.CollectMacOSStat(n, m)

	for {
		s := models.Stat{}
		select {
		case s = <-statChan:
			fmt.Printf(s.String())
		}
	}

	wg.Wait()
}

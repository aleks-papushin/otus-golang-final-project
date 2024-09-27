package collector

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aleks-papushin/system-monitor/internal/models"
)

func CollectMacOSStat(statChan chan *models.Stat) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	n := 5
	m := 15

	c := &MacOSStatCollector{
		Executor: &RealCommandExecutor{},
	}

	wg.Add(1)
	StartStatCollecting(n, m, statChan, &wg, c)

	wg.Wait()
}

func StartStatCollecting(n, m int, statChan chan *models.Stat, wg *sync.WaitGroup, c StatCollector) {
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

	//go func() {
	//	var stat *models.Stat
	//	for {
	//		select {
	//		case stat = <-statChan:
	//			fmt.Println("Received stat:", stat)
	//		}
	//	}
	//}()
}

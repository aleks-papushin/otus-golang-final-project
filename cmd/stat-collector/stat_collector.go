package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/aleks-papushin/system-monitor/internal/collector"
)

func main() {
	n := flag.Int("n", 0, "Output stat interval")
	m := flag.Int("m", 0, "Collecting stat interval")
	flag.Parse()

	if *n == 0 || *m == 0 {
		fmt.Println("Both arguments should be greater than 0")
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	c := collector.GetMacOSStatCollector()

	statChan := c.CollectStat(*n, *m)

	go func() {
		for {
			s, ok := <-statChan
			if !ok {
				break
			}
			fmt.Printf("%s\n", s.String())
		}
		wg.Done()
	}()

	wg.Wait()
}

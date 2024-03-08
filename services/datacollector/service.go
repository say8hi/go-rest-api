package main

import (
	"log"
	"service-datacollector/utils"
	"time"
)

func main() {
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				utils.FetchData()
			}
		}
	}()

	log.Println("Service started.")
	select {}
}

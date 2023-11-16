package main

import (
	"fmt"

	"github.com/robfig/cron"
)

func main() {
	scheduler := cron.New()

	scheduler.AddFunc("@every 10s", func() {
		fmt.Println("Executed every 10 seconds")
	})

	scheduler.Start()

	select {}
}
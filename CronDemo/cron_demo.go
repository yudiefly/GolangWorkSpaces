package main

import (
	"log"
	"time"

	"github.com/robfig/cron"
)

func main() {
	log.Println("Starting……")

	c := cron.New()
	c.AddFunc("* * * * * *", func() {
		log.Println("Run Clean Tags ……")
		time.Sleep(10 * time.Second)
		log.Println("Tag Clean End.")
	})

	c.AddFunc("* * * * * *", func() {
		log.Println("Run Clean Articles ……")
		time.Sleep(10 * time.Second)
		log.Println("Articles Clean End.")
	})

	c.Start()
	t1 := time.NewTimer(time.Second * 30)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 30)
		}
	}

}

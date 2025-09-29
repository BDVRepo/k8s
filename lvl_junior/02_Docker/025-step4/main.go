package main

import (
	"flag"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	url := flag.String("url", "https://google.com", "URL to be called")
	interval := flag.Int("interval", 1, "Interval between requests")
	flag.Parse()
	log.Infof("Started with arguments: %s, %d", url, interval)
	for {
		resp, err := http.Get(*url)
		if err != nil {
			log.Error("Error while http.Get: ", err.Error())
			time.Sleep(time.Duration(*interval) * time.Second)
			continue
		}

		log.Info(resp)
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}

package main

import (
	"fmt"
	"net/http"
	"time"

	"flag"

	"github.com/heptiolabs/healthcheck"
	log "github.com/sirupsen/logrus"
)

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path[1:] {
	case "hi":
		log.Info("path: /hi")
		fmt.Fprintf(w, "Hello world")
	}
}

func main() {

	serverURL := flag.String("url", "http://server-svc:8080/kubernetes", "URL to be called")
	waitSeconds := flag.Int("wait", 600, "Time seconds before creating 100 goroutines")
	flag.Parse()
	log.Infof("Started with arguments: serverURL=%v, waitSeconds=%v", *serverURL, *waitSeconds)

	health := healthcheck.NewHandler()
	health.AddReadinessCheck(
		"http-get-check",
		healthcheck.HTTPGetCheck(*serverURL, 100*time.Millisecond))
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))

	http.HandleFunc("/", handler)

	go http.ListenAndServe(":8086", health)
	go goroutineIncrementor(*waitSeconds)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func goroutineIncrementor(waitSeconds int) {
	time.Sleep(time.Duration(waitSeconds) * time.Second)
	for i := 0; i <= 100; i++ {
		go func() {
			time.Sleep(60 * time.Second)
			log.Info("Stop goroutine")
		}()
	}
}

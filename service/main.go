package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

func main() {
	delayStr := flag.String("repeat", "0", "Repeat request after interval")
	url := flag.String("url", "", "Url where to send requests")
	key := flag.String("key", "", "Secret key to use")
	flag.Parse()

	delay, err := time.ParseDuration(*delayStr)
	if err != nil {
		panic(err)
	}

	sendRequest(*url, *key)

	if delay == 0 {
		return
	}

	for _ = range time.Tick(delay) {
		sendRequest(*url, *key)
		runtime.GC()
	}
}

func sendRequest(url string, key string) {
	http.Get(fmt.Sprintf("%s?key=%s&uptime=%d", url, key, GetUptime()))
}

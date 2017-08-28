package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var param = map[string]*string{
	"repeat": nil,
	"url":    nil,
	"key":    nil,
}

func main() {
	param["repeat"] = flag.String("repeat", "0", "Repeat request after interval")
	param["url"] = flag.String("url", "", "Url where to send requests")
	param["key"] = flag.String("key", "", "Secret key to use")
	flag.Parse()

	// Environment variables can override cmdline parameter
	for p := range param {
		if os.Getenv(strings.ToUpper(p)) != "" {
			*param[p] = os.Getenv(strings.ToUpper(p))
		}
	}

	delay, err := time.ParseDuration(*param["repeat"])
	if err != nil {
		panic(err)
	}

	// Imidiadly send first request
	sendRequest(*param["url"], *param["key"])

	// Do not repeat
	if delay == 0 {
		return
	}

	// Repeat request forever
	for _ = range time.Tick(delay) {
		sendRequest(*param["url"], *param["key"])
		runtime.GC()
	}
}

func sendRequest(url string, key string) {
	fmt.Printf("%s?key=%s&uptime=%d\n", url, key, GetUptime())
	http.Get(fmt.Sprintf("%s?key=%s&uptime=%d", url, key, GetUptime()))
}

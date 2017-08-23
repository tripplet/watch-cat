package main

import (
	"time"
	"flag"
	"fmt"
	"net/http"
	"runtime"
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

	ticker := time.Tick(delay)

	SendRequest(*url, *key)
	
	if delay == 0 {
		return
	} 

	for now := range ticker {
		fmt.Printf("%v %d\n", now, GetUptime())
		SendRequest(*url, *key)
		runtime.GC()
	}
}

func SendRequest(url string, key string) {
	http.Get(fmt.Sprintf("%s?key=%s&uptime=%d", url, key, GetUptime()))
}
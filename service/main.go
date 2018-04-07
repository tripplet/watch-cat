package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type params struct {
	repeat   string
	url      string
	key      string
	nouptime bool
	verbose  bool
}

var param params

func main() {
	flag.StringVar(&param.repeat, "repeat", "0", "Repeat request after interval, valid units are 'ms', 's', 'm', 'h' e.g. 2m30s")
	flag.StringVar(&param.url, "url", "", "Url where to send requests")
	flag.StringVar(&param.key, "key", "", "Secret key to use")
	flag.BoolVar(&param.nouptime, "nouptime", false, "Do not send uptime in heartbeat requests")
	flag.BoolVar(&param.verbose, "verbose", false, "Verbose mode")
	flag.Parse()

	// Try to 	use environment variables if parameter is not set via cmdline
	log("Parameter:")
	cmd := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { cmd[f.Name] = true })

	paramStruct := reflect.TypeOf(param)
	for idx := 0; idx < paramStruct.NumField(); idx++ {
		name := paramStruct.Field(idx).Name
		env := os.Getenv(strings.ToUpper(name))
		_, isCmdLineParameter := cmd[name]

		if env != "" && !isCmdLineParameter {
			flag.Set(name, env)
		}

		log(name, "=", reflect.ValueOf(param).Field(idx))
	}

	delay, err := time.ParseDuration(param.repeat)
	if err != nil {
		panic(err)
	}

	// Immediately send first heartbeat
	sendRequest()

	// Do not repeat
	if delay <= 0 {
		return
	}

	// Repeat heartbeat forever
	for _ = range time.Tick(delay) {
		sendRequest()
		runtime.GC()
	}
}

func sendRequest() {
	var url string

	if param.nouptime {
		url = fmt.Sprintf("%s?key=%s", param.url, param.key)
	} else {
		url = fmt.Sprintf("%s?key=%s&uptime=%d", param.url, param.key, GetUptime())
	}

	log()
	log("Sending:", url)
	resp, err := http.Get(url)

	if err != nil {
		log(">>", err)
	} else {
		log(">>", resp.Status)
		defer resp.Body.Close()
	}
}

func log(l ...interface{}) {
	if param.verbose {
		fmt.Println(l...)
	}
}

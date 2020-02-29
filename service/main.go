package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

type params struct {
	repeat   string
	url      string
	key      string
	nouptime bool
	verbose  bool
	checkdns int
	timeout  int
}

var config params

func main() {
	// OS specific preparations
	OSSpecificPrepare()

	parseParameter()

	// Perform OS specific actions after paramter parsing has been done
	OSSpecific()

	delay, err := time.ParseDuration(config.repeat)
	if err != nil {
		panic(err)
	}

	if config.checkdns > 0 {
		checkDNS()
	}

	client := &http.Client{
		Timeout: time.Second * time.Duration(config.timeout),
	}

	// Immediately send first heartbeat
	sendRequest(client)

	// Do not repeat
	if delay <= 0 {
		os.Exit(0)
	}

	// Make the force garbage collection very aggressive
	debug.SetGCPercent(1)

	// Repeat heartbeat forever
	for range time.Tick(delay) {
		go sendRequestAndCleanup(client)
	}
}

func sendRequest(client *http.Client) {
	url := config.url

	if !config.nouptime || config.key != "" {
		url = url + "?"
	}

	if config.key != "" {
		url = url + fmt.Sprintf("key=%s", config.key)
	}

	if !config.nouptime {
		if config.key != "" {
			url = url + "&"
		}

		url = url + fmt.Sprintf("uptime=%d", GetUptime())
	}

	log()
	log("- Sending:", url)

	resp, err := client.Get(url)

	if err != nil {
		log("  >>", err)
	} else {
		defer resp.Body.Close()
		log("  >>", resp.Status)
	}
}

func sendRequestAndCleanup(client *http.Client) {
	sendRequest(client)
	runtime.GC()
}

func log(l ...interface{}) {
	if config.verbose {
		fmt.Println(l...)
	}
}

func parseParameter() {
	flag.StringVar(&config.repeat, "repeat", "0", "Repeat request after interval, valid units are 'ms', 's', 'm', 'h' e.g. 2m30s")
	flag.StringVar(&config.url, "url", "", "Url where to send requests")
	flag.StringVar(&config.key, "key", "", "Secret key to use")
	flag.IntVar(&config.timeout, "timeout", 60, "Timeout for http request in seconds")
	flag.IntVar(&config.checkdns, "checkdns", 0, "Check dns every x seconds before first request, for faster inital signal in case of long allowed timeout")
	flag.BoolVar(&config.nouptime, "nouptime", false, "Do not send uptime in heartbeat requests")
	flag.BoolVar(&config.verbose, "verbose", false, "Verbose mode")
	flag.Parse()

	// Try to use environment variables if parameter is not set via cmdline
	log("- Parameter:")
	cmd := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { cmd[f.Name] = true })

	paramStruct := reflect.TypeOf(config)
	for idx := 0; idx < paramStruct.NumField(); idx++ {
		name := paramStruct.Field(idx).Name
		env := os.Getenv(strings.ToUpper(name))
		_, isCmdLineParameter := cmd[name]

		if env != "" && !isCmdLineParameter {
			flag.Set(name, env)
		}

		log(" ", name, "=", reflect.ValueOf(config).Field(idx))
	}

	log()

	if config.url == "" {
		fmt.Println("No url provided use \"-help\" to see all supported command line arguments")
		os.Exit(-1)
	}
}

func checkDNS() {
	log("- Checking for DNS...")
	url, err := url.Parse(config.url)
	if err != nil {
		panic(err)
	}

	host := url.Hostname()
	resolve := &net.Resolver{}

	log("  Trying to resolve: " + host)
	idx := 0
	for {
		log("  Try:", idx)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.checkdns)*time.Second)
		defer cancel()
		_, err := resolve.LookupHost(ctx, host)

		if err == nil {
			log("  DNS available")
			break
		}

		idx++
		log()
	}
}

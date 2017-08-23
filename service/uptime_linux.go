package main

import (
	"syscall"
)

func GetUptime() int {
	sysinfo := syscall.Sysinfo_t{}

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return -1
	}

	return int(sysinfo.Uptime)
}
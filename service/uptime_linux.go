package main

import (
	"syscall"
)

// OSSpecificPrepare for system specific preparations
func OSSpecificPrepare() {}

// OSSpecific to perform system specific actions after paramter parsing has been done
func OSSpecific() {}

// GetUptime returns the system uptime
// See: https://github.com/cloudfoundry/gosigar (Apache 2 license)
func GetUptime() int {
	sysinfo := syscall.Sysinfo_t{}

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return -1
	}

	return int(sysinfo.Uptime)
}

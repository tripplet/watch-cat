package main

import (
	"syscall"
	"time"
)

var (
	kernel32Dll        = syscall.MustLoadDLL("kernel32")
	procGetTickCount64 = kernel32Dll.MustFindProc("GetTickCount64").Addr()
)

// GetUptime returns the system uptime
// See: https://github.com/cloudfoundry/gosigar (Apache 2 license)
func GetUptime() int {
	count, _, err := syscall.Syscall(procGetTickCount64, 0, 0, 0, 0)
	if err != 0 {
		return -1
	}
	return int((time.Duration(count) * time.Millisecond).Seconds())
}

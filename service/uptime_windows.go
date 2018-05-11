package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/windows/svc/mgr"
)

var installService = false

var (
	kernel32Dll        = syscall.MustLoadDLL("kernel32")
	procGetTickCount64 = kernel32Dll.MustFindProc("GetTickCount64").Addr()
	*installService  = false
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

// SystemSpecificPrepare for system specific preparations
func OSSpecificPrepare() {
	flag.BoolVar(&installService, "install", false, "Install the service")
}

// SystemSpecific to perform system specific actions after paramter parsing has been done
func OSSpecific() {
	if installService {
		var args []string
		flag.Visit(func(f *flag.Flag) {
			if f.Name != "install" {
				args = append(args, "--"+f.Name, f.Value.String())
			}
		})

		fmt.Println(args)
		//InstallService("watchcat", args)
		os.Exit(0)
	}
}

func InstallService(serviceName string, args string[]) error {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}

	serviceMgr, err := mgr.Connect()
	if err != nil {
		panic(err)
	}
	defer serviceMgr.Disconnect()

	// Check if service already installed
	service, err := serviceMgr.OpenService(serviceName)
	if err == nil {
		panic(err)
	}
	defer service.Close()

	service, err = serviceMgr.CreateService(serviceName, exe, mgr.Config{}, args...)
	if err != nil {
		return err
	}
	defer s.Close()

	err = s.Start("is", "manual-started")
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}

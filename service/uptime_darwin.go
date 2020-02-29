package main

import (
	"bytes"
	"encoding/binary"
	"syscall"
	"time"
	"unsafe"
)

// OSSpecificPrepare for system specific preparations
func OSSpecificPrepare() {}

// OSSpecific to perform system specific actions after paramter parsing has been done
func OSSpecific() {}

// GetUptime returns the system uptime
// See: https://github.com/cloudfoundry/gosigar (Apache 2 license)
func GetUptime() int {
	tv := syscall.Timeval32{}

	if err := sysctlbyname("kern.boottime", &tv); err != nil {
		return -1
	}

	return int(time.Since(time.Unix(int64(tv.Sec), int64(tv.Usec)*1000)).Seconds())
}

// generic Sysctl buffer unmarshalling
func sysctlbyname(name string, data interface{}) (err error) {
	val, err := syscall.Sysctl(name)
	if err != nil {
		return err
	}

	buf := []byte(val)

	switch v := data.(type) {
	case *uint64:
		*v = *(*uint64)(unsafe.Pointer(&buf[0]))
		return
	}

	bbuf := bytes.NewBuffer([]byte(val))
	return binary.Read(bbuf, binary.LittleEndian, data)
}

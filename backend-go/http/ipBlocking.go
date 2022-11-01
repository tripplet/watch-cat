package http

import (
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type blockedIP struct {
	InvalidRequests uint32
	LastInvalid     time.Time
	BlockedUntil    time.Time
}

var blockedIPs = make(map[string]*blockedIP)

func middlewareIPBlocking() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			return
		}

		entry, exists := blockedIPs[ip]
		if !exists {
			// IP no blocked, continue processing the request
			return
		}

		if (entry.BlockedUntil != time.Time{}) && time.Now().UTC().Before(entry.BlockedUntil) {
			c.AbortWithStatus(http.StatusForbidden)
		} else {
			// Unblock the IP
			delete(blockedIPs, ip)
		}
	}
}

func incBlockIP(c *gin.Context) {
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return
	}

	entry, exists := blockedIPs[ip]

	if !exists {
		entry = &blockedIP{
			InvalidRequests: 0,
			LastInvalid:     time.Time{},
			BlockedUntil:    time.Time{},
		}
		blockedIPs[ip] = entry
	}

	atomic.AddUint32(&entry.InvalidRequests, 1)
	entry.LastInvalid = time.Now().UTC()

	if entry.InvalidRequests >= 9 {
		entry.BlockedUntil = time.Now().UTC().Add(1 * time.Hour)
	}
}

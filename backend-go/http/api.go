package http

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"watchcat/models"
	"watchcat/taskQueue"

	"github.com/gin-gonic/gin"
)

// Get the secret from the request
func extractSecretFromRequest(c *gin.Context) string {
	authHeaderParts := strings.Split(c.GetHeader("Authorization"), " ")

	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		return ""
	}

	return authHeaderParts[1]
}

func (env *HttpEnv) jobUpdate(c *gin.Context) {
	ctx := c.Request.Context()

	secret := extractSecretFromRequest(c)
	if secret == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	job, err := models.GetWatchJobForSecret(env.global, ctx, secret)
	if err != nil {
		log.Println(err)
		return
	} else if job == nil {
		incBlockIP(c)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	log.Printf("Job %s updated", job.Name)

	// Found the job now update LastSeen and perform the necessary actions
	job.LastSeen = time.Now().UTC()

	remoteIP, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		log.Println(err)
		return
	}

	if remoteIP != job.LastIP && remoteIP != job.LastIPv4 && remoteIP != job.LastIPv6 {
		models.CreateLogEntry(env.global, job.ID, "IP Change", fmt.Sprintf("IP changed - new IP: %s", remoteIP))
	}

	// Update the IP
	job.LastIP = remoteIP
	if strings.Contains(remoteIP, ":") {
		job.LastIPv6 = remoteIP
	} else {
		job.LastIPv4 = remoteIP
	}

	// Update the reported uptime
	uptime := time.Duration(0)
	uptimeStr := c.Query("uptime")
	if uptimeStr != "" {
		uptimeSeconds, err := strconv.ParseInt(uptimeStr, 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			log.Println(err)
			return
		}

		uptime = time.Duration(uptimeSeconds) * time.Second

		// Check for reboot
		if job.Uptime != 0 && job.Uptime > uptime {
			models.CreateLogEntry(env.global, job.ID, "Reboot", fmt.Sprintf("Reboot detected - old uptime: %s", job.Uptime))

			// Perform the reboot actions
			for _, action := range job.RebootActions {
				go action.Run()
			}
		}

		job.Uptime = uptime
	}

	// job got back online
	if job.Status == "offline" {
		job.Status = "online"
		models.CreateLogEntry(env.global, job.ID, "Back Online", fmt.Sprintf("%s is back online - IP: %s", job.Name, remoteIP))

		// Perform all back online actions
		for _, action := range job.BackOnlineActions {
			go action.Run()
		}
	}

	// Delete previous (waiting) task
	if job.TaskId != 0 {
		if err := env.global.Dispatcher.Cancel(job.TaskId); err != nil {
			log.Printf("Error deleting task %d: %s", job.TaskId, err)
		}
	}

	//newTaskName := job.Name + "_" + time.Now().UTC().Format("2006-01-02_15-04-05")

	// Schedule new task
	newTaskId, err := env.global.Dispatcher.Schedule(taskQueue.Task{
		StartIn: time.Second * time.Duration(job.Interval),
		Fn:      func() { fmt.Println("test") },
	})

	if err != nil {
		log.Printf("Error scheduling task %d: %s", newTaskId, err)
		return
	}

	job.TaskId = newTaskId
	env.global.Database.Save(&job)
}

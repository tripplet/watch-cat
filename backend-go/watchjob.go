package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type watchJob struct {
	ID                uint
	Name              string
	LastSeen          time.Time
	Uptime            time.Duration
	Enabled           bool
	Status            string
	Interval          int
	Secret            string
	LastIP            string
	LastIPv4          string
	LastIPv6          string
	TimeoutActions    []actionData `gorm:"many2many:timeout_actions;"`
	BackOnlineActions []actionData `gorm:"many2many:backonline_actions;"`
	RebootActions     []actionData `gorm:"many2many:reboot_actions;"`
}

// getWatchJobForSecret gets a watchJob from the firestore api based on the given secret
func getWatchJobForSecret(ctx context.Context, secret string) (*watchJob, error) {
	if secret == "" {
		return nil, nil
	}

	var job watchJob
	result := db.Where("secret = ?", secret).First(&job)
	if result.Error != nil {
		return nil, result.Error
	}

	return &job, nil
}

// Get the secret from the request
func getSecretFromRequest(c *gin.Context) string {
	authHeaderParts := strings.Split(c.GetHeader("Authorization"), " ")

	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		return ""
	}

	return authHeaderParts[1]
}

func jobUpdate(c *gin.Context) {
	ctx := c.Request.Context()

	secret := getSecretFromRequest(c)
	if secret == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	job, err := getWatchJobForSecret(ctx, secret)
	if err != nil {
		log.Println(err)
		return
	} else if job == nil {
		incBlockIP(c)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	log.Printf("Job %s updated", job.Name)

	// Found the job now update it and perform the necessary actions
	job.LastSeen = time.Now().UTC()

	remoteIP, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		log.Println(err)
		return
	}

	if remoteIP != job.LastIP && remoteIP != job.LastIPv4 && remoteIP != job.LastIPv6 {
		// TODO LogEntry.log_event(self.key(), "Info', 'IP changed - new IP: ' + remote_ip)
	}

	job.LastIP = remoteIP
	if strings.Contains(remoteIP, ":") {
		job.LastIPv6 = remoteIP
	} else {
		job.LastIPv4 = remoteIP
	}

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

		if job.Uptime != 0 && job.Uptime > uptime {
			//TODO LogEntry.log_event(self.key(), 'Reboot', 'Reboot - Previous uptime: ' + str(timedelta(seconds=self.uptime)))

			//for _, _ := range job.RebootActions {
			// TODO
			//}

		}

		job.Uptime = uptime
	}

	// job got back online
	if job.Status == "offline" {
		job.Status = "online"
		// TODOLogEntry.log_event(self.key(), "Info", "Job back online - IP: " + remote_ip)

		// Perform all back online actions
		//for _, _ := range job.BackOnlineActions {
		// TODO
		//}
	}

	// Delete previous (waiting) task
	// if job.TaskName != "" {
	/*taskqueue.Delete(c.Request.Context(), &taskqueue.Task{Name: job.TaskName}, "")
	if err != nil {
		log.Println(err)
		return
	}
	*/
	// }

	newTaskName := job.Name + "_" + time.Now().UTC().Format("2006-01-02_15-04-05")
	fmt.Println(newTaskName)

	// appEngineClient, err := cloudtasks.NewClient(ctx, option.WithCredentialsFile("appEngineAccount.json"))

	// req := &taskspb.CreateQueueRequest{
	// 	// TODO: Fill request struct fields.
	// }

	// resp, err := c.CreateQueue(ctx, req)
	// if err != nil {
	// 	// TODO: Handle error.
	// }

	/*// Create task to be executed in updated no called in interval minutes
	newTask, err := taskqueue.Add(gae,
		&taskqueue.Task{
			Name:    newTaskName,
			Delay:   time.Duration(job.Interval+2) * time.Minute,
			Path:    "/task",
			Method:  "POST",
			Payload: []byte(jobDoc.Path),
		}, "timeouts")

	//(name=task_name, url="/task", params={"key": self.key()}, )

	_ = newTask
	job.TaskName = newTaskName
	*/

	// jobDoc.Set(ctx, job)
}

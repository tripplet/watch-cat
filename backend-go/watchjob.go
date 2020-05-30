package main

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

type watchJob struct {
	Name              string        `json:"name"`
	LastSeen          time.Time     `json:"last_seen"`
	Uptime            time.Duration `json:"uptime"`
	Enabled           bool          `json:"enabled"`
	Status            string        `json:"status"`
	Interval          int           `json:"interval"`
	Secret            string        `json:"secret"`
	LastIP            string        `json:"last_ip"`
	LastIPv4          string        `json:"last_ipv4"`
	LastIPv6          string        `json:"last_ipv6"`
	TaskName          string        `json:"task_name"`
	TimeoutActions    []int64       `json:"actions_timeout"`
	BackOnlineActions []int64       `json:"actions_back_online"`
	RebootActions     []int64       `json:"actions_reboot"`
}

// getWatchJobForSecret gets a watchJob from the firestore api based on the given secret
func getWatchJobForSecret(ctx context.Context, secret string) (*watchJob, *firestore.DocumentRef, error) {
	if secret == "" {
		return nil, nil, nil
	}

	return nil, nil, nil

	// jobDoc, err := client.Collection("WatchJob").Where("secret", "==", secret).Documents(ctx).GetAll()
	// if err != nil {
	// 	return nil, nil, err
	// }

	// if len(jobDoc) != 1 {
	// 	return nil, nil, nil
	// }

	// var job watchJob
	// if err := jobDoc[0].DataTo(&job); err != nil {
	// 	return nil, nil, err
	// }

	// return &job, jobDoc[0].Ref, nil
}

func jobUpdate(c *gin.Context) {
	// ctx := c.Request.Context()

	// job, jobDoc, err := getWatchJobForSecret(ctx, c.Param("key"))
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// } else if job == nil {
	// 	incBlockIP(client, c)
	// 	c.AbortWithStatus(http.StatusNotFound)
	// 	return
	// }

	// job.LastSeen = time.Now().UTC()

	// remoteIP, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// if remoteIP != job.LastIP && remoteIP != job.LastIPv4 && remoteIP != job.LastIPv6 {
	// 	// TODO LogEntry.log_event(self.key(), "Info', 'IP changed - new IP: ' + remote_ip)
	// }

	// job.LastIP = remoteIP
	// if strings.Contains(remoteIP, ":") {
	// 	job.LastIPv6 = remoteIP
	// } else {
	// 	job.LastIPv4 = remoteIP
	// }

	// uptime := time.Duration(0)
	// uptimeStr := c.Query("uptime")
	// if uptimeStr != "" {
	// 	uptimeSeconds, err := strconv.ParseInt(uptimeStr, 10, 64)
	// 	if err != nil {
	// 		incBlockIP(client, c)
	// 		c.AbortWithStatus(http.StatusBadRequest)
	// 		log.Println(err)
	// 		return
	// 	}

	// 	uptime = time.Duration(uptimeSeconds) * time.Second

	// 	if job.Uptime != 0 && job.Uptime > uptime {
	// 		//TODO LogEntry.log_event(self.key(), 'Reboot', 'Reboot - Previous uptime: ' + str(timedelta(seconds=self.uptime)))

	// 		//for _, _ := range job.RebootActions {
	// 		// TODO
	// 		//}

	// 	}

	// 	job.Uptime = uptime
	// }

	// // job got back online
	// if job.Status == "offline" {
	// 	job.Status = "online"
	// 	// TODOLogEntry.log_event(self.key(), "Info", "Job back online - IP: " + remote_ip)

	// 	// Perform all back online actions
	// 	//for _, _ := range job.BackOnlineActions {
	// 	// TODO
	// 	//}
	// }

	// // Delete previous (waiting) task
	// if job.TaskName != "" {
	// 	/*taskqueue.Delete(c.Request.Context(), &taskqueue.Task{Name: job.TaskName}, "")
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	*/
	// }

	// newTaskName := job.Name + "_" + time.Now().UTC().Format("2006-01-02_15-04-05")
	// fmt.Println(newTaskName)

	// // appEngineClient, err := cloudtasks.NewClient(ctx, option.WithCredentialsFile("appEngineAccount.json"))

	// // req := &taskspb.CreateQueueRequest{
	// // 	// TODO: Fill request struct fields.
	// // }

	// // resp, err := c.CreateQueue(ctx, req)
	// // if err != nil {
	// // 	// TODO: Handle error.
	// // }

	// /*// Create task to be executed in updated no called in interval minutes
	// newTask, err := taskqueue.Add(gae,
	// 	&taskqueue.Task{
	// 		Name:    newTaskName,
	// 		Delay:   time.Duration(job.Interval+2) * time.Minute,
	// 		Path:    "/task",
	// 		Method:  "POST",
	// 		Payload: []byte(jobDoc.Path),
	// 	}, "timeouts")

	// //(name=task_name, url="/task", params={"key": self.key()}, )

	// _ = newTask
	// job.TaskName = newTaskName
	// */

	// jobDoc.Set(ctx, job)
}

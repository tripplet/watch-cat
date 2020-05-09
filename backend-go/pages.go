package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func getJobs(ctx context.Context) ([]watchJob, error) {
	jobs, err := client.Collection("WatchJob").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var renderJobs []watchJob

	for _, jobDoc := range jobs {
		var job watchJob

		if err := jobDoc.DataTo(&job); err != nil {
			return nil, err
		}

		renderJobs = append(renderJobs, job)
	}

	return renderJobs, nil
}

func handleRootPage(c *gin.Context) {
	jobs, err := getJobs(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	c.HTML(http.StatusOK, "root.htm", gin.H{"jobs": jobs})
}

func handleLogPage(c *gin.Context) {
	//c.HTML(http.StatusOK, "log.htm", gin.H{"job": jobs[0]})
}

func handleDebugPage(c *gin.Context) {
	jobs, err := getJobs(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var sb strings.Builder
	sb.WriteString("<b><i>ServerTime: </i></b>")

	sb.WriteString(time.Now().UTC().In(tzLocation).Format("15:04:05 - 02.01.2006"))
	sb.WriteString("<br><br>")

	for _, job := range jobs {
		sb.WriteString(fmt.Sprintf("<b>%s</b> <a href=\"/notify/%s\">[testNotification]</a><br>%s<br>%s<br>%s<br><br>",
			job.Name,
			job.Name,
			formatDatetime(job.LastSeen),
			job.LastIP,
			formatTimespan(job.Uptime)))
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(sb.String()))
}

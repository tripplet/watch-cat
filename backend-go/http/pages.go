package http

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strings"
	"text/template"
	"time"
	"watchcat/env"
	"watchcat/models"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var tzLocation *time.Location

type HttpEnv struct {
	global *env.Env
}

func Startup(globalEnv *env.Env) {
	var env HttpEnv = HttpEnv{globalEnv}

	router := gin.Default()

	// Load timezone TODO from env
	tzLocation, _ = time.LoadLocation(viper.GetString("web.timezone"))

	router.SetTrustedProxies(viper.GetStringSlice("http.trusted_proxies"))

	// Add template functions
	router.SetFuncMap(template.FuncMap{
		"format_datetime": formatDatetime,
		"add_breakchars":  formatAddBreakChars,
		"format_timespan": formatTimespan,
		"lower":           formatToLower,
	})

	// Load templates
	router.LoadHTMLGlob("templates/*.htm")

	// Static assets
	mime.AddExtensionType(".js", "application/javascript")
	router.Static("/assets", "./static/")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	router.StaticFile("/robots.txt", "./static/robots.txt")

	// Routes
	router.GET("/", env.handleRootPage)
	router.GET("/log/:job", env.handleLogPage)
	router.GET("/job/:job", env.handleJobPage)

	router.GET("/debug", env.handleDebugPage)
	router.POST("/create/:job", env.create)

	router.GET("/notify/:job", env.notifyTest)

	router.Use(middlewareIPBlocking()) // Add IP blocking middleware

	// Public API
	api := router.Group("/api/v2")
	api.POST("/job/update", env.jobUpdate)

	router.Run(viper.GetString("http.listen"))
}

func (env *HttpEnv) handleRootPage(c *gin.Context) {
	jobs, err := models.GetJobs(env.global, c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	c.HTML(http.StatusOK, "root.htm", gin.H{"jobs": jobs})
}

func (env *HttpEnv) handleLogPage(c *gin.Context) {
	//c.HTML(http.StatusOK, "log.htm", gin.H{"job": jobs[0]})
}

func (env *HttpEnv) handleJobPage(c *gin.Context) {
	jobName := c.Param("job")
	if jobName == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	job, err := models.GetJob(env.global, jobName)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		log.Println(err)
		return
	}

	c.HTML(http.StatusOK, "job.htm", gin.H{"job": job})
}

func (env *HttpEnv) handleDebugPage(c *gin.Context) {
	jobs, err := models.GetJobs(env.global, c.Request.Context())
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

func (env *HttpEnv) create(c *gin.Context) {
	newJobName := c.Param("job")
	c.String(http.StatusOK, newJobName)

	var key string
	for {
		randomBytes := make([]byte, 128)
		_, err := rand.Read(randomBytes)
		if err != nil {
			panic(err)
		}

		key = base64.StdEncoding.EncodeToString(randomBytes)
		key = strings.Replace(key, "=", "", -1)
		key = strings.Replace(key, "+", "", -1)
		key = strings.Replace(key, "/", "", -1)

		if len(key) < 48 {
			continue
		} else {
			key = key[:48]
			break
		}
	}

	env.global.Database.Create(&models.WatchJob{
		Name:     newJobName,
		Enabled:  true,
		Interval: 5 * 60, // 5 minutes
		Secret:   key,
		Status:   "offline",
	})
}

func (env *HttpEnv) notifyTest(c *gin.Context) {
	//jobName := c.Param("job")

	// jobDoc, err := client.Collection("WatchJob").Where("name", "==", jobName).Documents(c.Request.Context()).Next()
	// if err != nil {
	// 	if status.Code(err) == codes.NotFound {
	// 		//incBlockIP(client, c)
	// 		c.AbortWithStatus(http.StatusNotFound)
	// 		return
	// 	}

	// 	log.Println(err)
	// 	return
	// }

	// var job watchJob
	// if err := jobDoc.DataTo(&job); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// for _, actionRef := range job.TimeoutActions {
	// 	doc, err := actionRef.Get(c.Request.Context())
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}

	// 	// TODO handle all type of actions?
	// 	var action pushoverAction
	// 	if err := doc.DataTo(&action); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}

	// 	action.Run()
	// }
}

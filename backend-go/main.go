package main

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var tzLocation *time.Location
var db *gorm.DB

func main() {
	var err error

	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&watchJob{})

	// Load timezone TODO from env
	tzLocation, _ = time.LoadLocation("Europe/Berlin")

	r := gin.Default()

	r.Use()

	// Add template functions
	r.SetFuncMap(template.FuncMap{
		"format_datetime": formatDatetime,
		"add_breakchars":  formatAddBreakChars,
		"format_timespan": formatTimespan,
		"lower":           formatToLower,
	})

	// Load templates
	r.LoadHTMLGlob("templates/*.htm")

	// Static assets
	mime.AddExtensionType(".js", "application/javascript")
	r.Static("/assets", "./static/")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.StaticFile("/robots.txt", "./static/robots.txt")

	// Routes
	r.GET("/", handleRootPage)
	r.GET("/log/:job", handleLogPage)
	r.GET("/job/:job", handleJobPage)

	r.GET("/debug", handleDebugPage)
	r.GET("/create/:job", create)

	r.GET("/notify/:job", notifyTest)
	r.GET("/task/:key", executeTask)

	r.Use(middlewareIPBlocking()) // Add IP blocking middleware

	// Public API
	api := r.Group("/api/v2")
	api.POST("/job/update", jobUpdate)

	r.Run("127.0.0.1:8080")
}

func create(c *gin.Context) {
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

	db.Create(&watchJob{
		Name:     newJobName,
		Enabled:  true,
		Interval: 5 * 60, // 5 minutes
		Secret:   key,
		Status:   "offline",
	})
}

func notifyTest(c *gin.Context) {
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

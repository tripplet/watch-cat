package main

import (
	"html/template"
	"log"
	"mime"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var tzLocation *time.Location
var projectID string
var client *firestore.Client

func main() {
	projectID = "timo24-fd784" // os.Getenv("GOOGLE_CLOUD_PROJECT")
	// if projectID == "" {
	// 	log.Fatal("GOOGLE_CLOUD_PROJECT must be set")
	// }

	var err error
	ctx := context.Background()
	client, err = firestore.NewClient(ctx, projectID, option.WithCredentialsFile("serviceAccountKey.json"))
	// client, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}

	defer client.Close()

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

	r.GET("/debug", handleDebugPage)
	r.GET("/create/:job", create)

	r.GET("/notify/:job", notifyTest)
	r.GET("/task/:key", executeTask)

	// Public API
	api := r.Group("/api/v2")
	api.Use(middlewareIPBlocking(client)) // Add IP blocking middleware
	api.GET("/job/:key", jobUpdate)

	r.AppEngine = true

	r.Run("127.0.0.1:8080")
}

func create(c *gin.Context) {
	newJobName := c.Param("job")
	c.String(http.StatusOK, newJobName)
}

func notifyTest(c *gin.Context) {
	jobName := c.Param("job")

	if jobName == "" {
		incBlockIP(client, c)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	jobDoc, err := client.Collection("WatchJob").Where("name", "==", jobName).Documents(c.Request.Context()).Next()
	if err != nil {
		if status.Code(err) == codes.NotFound {
			incBlockIP(client, c)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Println(err)
		return
	}

	var job watchJob
	if err := jobDoc.DataTo(&job); err != nil {
		log.Println(err)
		return
	}

	for _, actionRef := range job.TimeoutActions {
		doc, err := actionRef.Get(c.Request.Context())
		if err != nil {
			log.Println(err)
			return
		}

		// TODO handle all type of actions?
		var action pushoverAction
		if err := doc.DataTo(&action); err != nil {
			log.Println(err)
			return
		}

		action.Run()
	}
}

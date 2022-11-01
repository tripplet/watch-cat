package main

import (
	"fmt"
	"log"
	"runtime/debug"

	"watchcat/env"
	"watchcat/http"
	"watchcat/models"
	"watchcat/taskQueue"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Read the config file
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %s", err)
	}

	setConfigDefaults()

	// Read the build info
	info, ok := debug.ReadBuildInfo()
	if ok {
		fmt.Printf("Version: %s", getVersionInfo(info))
	}

	// Open the database connection
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.WatchJob{})

	// Create environment for dependency injection
	env := &env.Env{
		Database:   db,
		Dispatcher: taskQueue.CreateDispatcher(),
	}

	models.CleanupTasks(env)
	env.Dispatcher.Start()

	http.Startup(env)
}

func setConfigDefaults() {
	viper.SetDefault("web.timezone", "Europe/Berlin")
}

func getVersionInfo(info *debug.BuildInfo) string {
	revision := ""
	dirty := false
	buildTime := ""

	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			revision = setting.Value
		}

		if setting.Key == "vcs.modified" && setting.Value == "true" {
			dirty = true
		}

		if setting.Key == "vcs.time" {
			buildTime = setting.Value
		}
	}

	if len(revision) > 12 {
		revision = revision[:10]
	}

	if dirty {
		revision += " (modified)"
	}

	return fmt.Sprintf("%s %s", revision, buildTime)
}

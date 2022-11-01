package env

import (
	"watchcat/taskQueue"

	"gorm.io/gorm"
)

// Global environment
type Env struct {
	Database   *gorm.DB
	Dispatcher taskQueue.Dispatcher
}

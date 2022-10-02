package main

import (
	"log"
)

type Task struct {
}

type Dispatcher interface {
	// Schedule a new task on the dispatcher
	Schedule(task Task) (string, error)

	// Cancel a task
	Cancel(taskName string) error

	// Starts the dispatcher
	Start()
}

type TaskDispatcher struct {
}

func CreateDispatcher() Dispatcher {
	return &TaskDispatcher{}
}

func (d *TaskDispatcher) Schedule(task Task) (string, error) {
	//log.Printf("scheduling Task %s: %s with deadline %s", taskName, payload, deadline)
	return "", nil
}

func (d *TaskDispatcher) Start() {

}

func (d *TaskDispatcher) Cancel(taskName string) error {
	log.Printf("deleting Task %s", taskName)
	return nil
}

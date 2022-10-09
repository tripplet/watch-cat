package taskQueue

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

type Task struct {
	StartIn time.Duration
	Fn      func()
}

type Dispatcher interface {
	// Schedule a new task on the dispatcher
	Schedule(task Task) (uint64, error)

	// Cancel a task
	Cancel(taskId uint64) error

	// Starts the dispatcher
	Start()
}

type TaskDispatcher struct {
	taskId uint64
	tasks  map[uint64]context.CancelFunc
}

func CreateDispatcher() Dispatcher {
	return &TaskDispatcher{taskId: 0, tasks: make(map[uint64]context.CancelFunc)}
}

func (d *TaskDispatcher) Schedule(task Task) (uint64, error) {
	// Get id for new task
	taskId := atomic.AddUint64(&d.taskId, 1)
	log.Printf("scheduling Task %d with deadline %s", taskId, task.StartIn)

	ctx, cancel := context.WithCancel(context.Background())
	d.tasks[taskId] = cancel

	go func() {
		defer cancel()
		defer d.Cancel(taskId)

		select {
		case <-time.After(task.StartIn):
			log.Printf("Executing task %d\n", taskId)
			task.Fn()
		case <-ctx.Done():
			log.Printf("Canceling task %d\n", taskId)
		}
	}()

	return taskId, nil
}

func (d *TaskDispatcher) Start() {}

func (d *TaskDispatcher) Cancel(taskId uint64) error {
	log.Printf("Deleting task %d", taskId)

	if cancelFunc, ok := d.tasks[taskId]; ok {
		delete(d.tasks, taskId)
		cancelFunc()
	} else {
		return fmt.Errorf("invalid task %d, could not find it", taskId)
	}

	return nil
}

package models

import (
	"context"
	"time"

	"watchcat/actions"
	"watchcat/env"
)

type WatchJob struct {
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
	TaskId            uint64
	TimeoutActions    []actions.ActionData `gorm:"many2many:timeout_actions;"`
	BackOnlineActions []actions.ActionData `gorm:"many2many:backonline_actions;"`
	RebootActions     []actions.ActionData `gorm:"many2many:reboot_actions;"`
}

func CleanupTasks(env *env.Env) error {
	jobs, err := GetJobs(env, context.Background())
	if err != nil {
		return err
	}

	// Reset the task id for all jobs
	for idx := range jobs {
		jobs[idx].TaskId = 0
	}

	env.Database.Save(jobs)
	return nil
}

// Gets a watchJob based on the given secret
func GetWatchJobForSecret(env *env.Env, ctx context.Context, secret string) (*WatchJob, error) {
	if secret == "" {
		return nil, nil
	}

	var job WatchJob
	result := env.Database.Where("secret = ?", secret).First(&job)
	if result.Error != nil {
		return nil, result.Error
	}

	return &job, nil
}

func GetJobs(env *env.Env, ctx context.Context) ([]WatchJob, error) {
	var jobs []WatchJob

	result := env.Database.Model(&WatchJob{}).Find(&jobs)
	if result.Error != nil {
		return nil, result.Error
	}

	return jobs, nil
}

func GetJob(env *env.Env, jobName string) (*WatchJob, error) {
	var job WatchJob
	result := env.Database.Where("name = ?", jobName).First(&job)

	if result.Error != nil {
		return nil, result.Error
	}

	return &job, nil
}

package cron

import "github.com/vinicius73/thecollector/pkg/tasks"

type Schedule struct {
	Action tasks.Action `yaml:"action"`
	Cron   []string     `yaml:"cron"`
}

type Schedules []Schedule

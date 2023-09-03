package main

import (
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/thecollector/app/config"
	"github.com/vinicius73/thecollector/pkg/cron"
	"github.com/vinicius73/thecollector/pkg/support"
)

var CronCmd = &cli.Command{
	Name:  "cron",
	Usage: "Start cron worker",
	Action: func(cmd *cli.Context) error {
		ctx, cancel := support.WithKillSignal(cmd.Context)

		defer cancel()

		cfg := *config.Ctx(ctx)

		worker, err := cron.New(cron.WorkerOptions{
			TargetDir:   cfg.TargetDir,
			Database:    cfg.Database,
			Timezone:    cfg.Timezone,
			Datasources: cfg.Datasources,
			SyncOptions: cfg.SyncOptions,
		})
		if err != nil {
			return err
		}

		return worker.Run(ctx, cfg.Schedules)
	},
}

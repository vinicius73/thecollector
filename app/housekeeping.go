package main

import (
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/thecollector/app/config"
	"github.com/vinicius73/thecollector/pkg/housekeeping"
	"github.com/vinicius73/thecollector/pkg/support"
	"github.com/vinicius73/thecollector/pkg/tasks"
)

var HousekeepingCmd = &cli.Command{
	Name:        "housekeeping",
	Description: "Housekeeping tasks: delete old dumps, etc.",
	Action: func(cmd *cli.Context) error {
		ctx, cancel := support.WithKillSignal(cmd.Context)

		defer cancel()

		cfg := *config.Ctx(ctx)

		return tasks.Housekeeping(ctx, housekeeping.Options{
			Config:  cfg.Housekeeping,
			BaseDir: cfg.TargetDir,
		})
	},
}

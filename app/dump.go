package main

import (
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/thecollector/app/config"
	"github.com/vinicius73/thecollector/pkg/support"
	"github.com/vinicius73/thecollector/pkg/tasks"
)

var DumpCmd = &cli.Command{
	Name:  "dump",
	Usage: "dump a database",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "dbname",
			Usage:    "database name",
			Required: true,
		},
	},
	Action: func(cmd *cli.Context) error {
		ctx, cancel := support.WithKillSignal(cmd.Context)

		defer cancel()

		cfg := *config.Ctx(ctx)

		_, err := tasks.DumbDB(ctx, tasks.DumpOptions{
			Database:  cfg.Database,
			TargetDir: cfg.TargetDir,
			Debug:     cfg.Logger.Level == "debug",
			DBName:    cmd.String("dbname"),
		})

		return err
	},
}

package main

import (
	"fmt"
	"os"
	"sort"

	zero "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/vinicius73/thecollector/app/config"
	"github.com/vinicius73/thecollector/pkg/logger"
	"github.com/vinicius73/thecollector/pkg/support"

	"github.com/vinicius73/thecollector/pkg/vars"
)

func main() {
	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "thecollector",
		Description:          "The Collector - Database Backup Tool",
		Usage:                "Database Backup Tool",
		Copyright:            "@vinicius73 - github.com/vinicius73/thecollector",
		Version:              vars.Version(),
		Commands:             []*cli.Command{DumpCmd, CronCmd},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "load configuration from",
				DefaultText: fmt.Sprintf("%s/thecollector.yml", support.GetBinDirPath()),
			},
			&cli.StringFlag{
				Name:        "level",
				Aliases:     []string{"l"},
				Usage:       "define log level",
				DefaultText: "info",
			},
		},
		Before: func(cmd *cli.Context) error {
			agentConfig, err := config.Load(cmd.String("config"))
			if err != nil {
				return err
			}

			logger.SetupLogger(agentConfig.Logger, cmd.String("level"), nil)

			log := logger.Logger("", nil)

			cmd.Context = log.WithContext(agentConfig.WithContext(cmd.Context))

			log.Debug().Msg("Config loaded")

			return nil
		},
	}

	cli.VersionPrinter = func(_ *cli.Context) {
		fmt.Println("The Collector - Database Backup Tool")
		fmt.Println("")
		fmt.Println(vars.VersionVerbose())
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		zero.Fatal().Err(err).Msg("fail run application")
	}
}

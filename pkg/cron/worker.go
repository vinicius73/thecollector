package cron

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
	"github.com/vinicius73/thecollector/pkg/cloud"
	"github.com/vinicius73/thecollector/pkg/database"
	"github.com/vinicius73/thecollector/pkg/tasks"
)

type WorkerOptions struct {
	Database    database.Config
	SyncOptions cloud.SyncOptions
	Datasources []string
	TargetDir   string
	Timezone    string
}

type Worker struct {
	options   WorkerOptions
	scheduler *gocron.Scheduler
}

func New(opt WorkerOptions) (*Worker, error) {
	w := &Worker{
		options: opt,
	}

	timezone, err := time.LoadLocation(opt.Timezone)
	if err != nil {
		return nil, err
	}

	w.scheduler = gocron.NewScheduler(timezone)

	return w, nil
}

func (w Worker) Run(ctx context.Context, entries Schedules) error {
	w.scheduler.Clear()
	w.scheduler.SingletonModeAll()

	logger := zerolog.Ctx(ctx).With().Str("worker", "cron").Logger()

	logger.Info().Msgf("registring entries (%s)", w.scheduler.Location().String())

	for _, entry := range entries {
		err := w.register(logger.WithContext(ctx), entry)
		if err != nil {
			logger.Error().Err(err).Msg("fail to registry cron entry")
			return err
		}
	}

	w.scheduler.StartAsync()

	logJobs := func() {
		for _, job := range w.scheduler.Jobs() {
			logger.Info().
				Strs("tags", job.Tags()).
				Msgf("next run: %s", job.NextRun().Format(time.RFC3339))
		}
	}

	go func() {
		logJobs()

		ticker := time.NewTicker(time.Minute * 45)

		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			logJobs()
		}
	}()

	<-ctx.Done()

	w.scheduler.Stop()

	return nil
}

func (w Worker) register(ctx context.Context, entry Schedule) error {
	logger := zerolog.Ctx(ctx).
		With().
		Str("action", string(entry.Action)).
		Logger()

	ctx = logger.WithContext(ctx)

	opt := tasks.Options{
		TargetDir:   w.options.TargetDir,
		Database:    w.options.Database,
		Datasources: w.options.Datasources,
		SyncOptions: w.options.SyncOptions,
		Debug:       false,
	}

	logger.Info().Msgf("registring cron entries")

	for _, cronExpression := range entry.Cron {
		logger.Info().Msgf("registring: %s", cronExpression)

		_, err := w.scheduler.
			Cron(cronExpression).
			Tag("action:"+string(entry.Action)).
			Do(entry.Action.Run, ctx, opt)
		if err != nil {
			return err
		}
	}

	return nil
}

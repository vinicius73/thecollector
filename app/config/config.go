package config

import (
	"context"

	"github.com/vinicius73/thecollector/pkg/cloud"
	"github.com/vinicius73/thecollector/pkg/cron"
	"github.com/vinicius73/thecollector/pkg/database"
	"github.com/vinicius73/thecollector/pkg/logger"
	"github.com/vinicius73/thecollector/pkg/vars"
)

type ctxKey struct{}

type App struct {
	Database    database.Config   `yaml:"database"`
	Logger      logger.Config     `yaml:"logger"`
	SyncOptions cloud.SyncOptions `yaml:"sync"`
	TargetDir   string            `yaml:"target_dir" default:"./dumps"`
	Timezone    string            `yaml:"timezone" default:"UTC"`
	Datasources []string          `yaml:"datasources"`
	Schedules   cron.Schedules    `yaml:"schedules"`
}

func Ctx(ctx context.Context) *App {
	cf, _ := ctx.Value(ctxKey{}).(*App)

	return cf
}

func (c App) Tags() map[string]interface{} {
	return map[string]interface{}{
		"v": vars.Version(),
	}
}

func (c *App) WithContext(ctx context.Context) context.Context {
	if cf, ok := ctx.Value(ctxKey{}).(*App); ok {
		if cf == c {
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, c)
}

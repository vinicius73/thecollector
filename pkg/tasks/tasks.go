package tasks

import (
	"context"

	"github.com/vinicius73/thecollector/pkg/cloud"
	"github.com/vinicius73/thecollector/pkg/database"
	"github.com/vinicius73/thecollector/pkg/errors"
	"github.com/vinicius73/thecollector/pkg/housekeeping"
)

type Action string

type Options struct {
	Debug        bool
	Database     database.Config
	SyncOptions  cloud.SyncOptions
	Housekeeping housekeeping.Config
	TargetDir    string
	Datasources  []string
}

const (
	ActionDump         Action = "dump"
	ActionHousekeeping Action = "housekeeping"
)

func (act Action) Run(ctx context.Context, opt Options) error {
	switch act {
	case ActionDump:
		return Dump(ctx, opt)
	case ActionHousekeeping:
		return Housekeeping(ctx, housekeeping.Options{
			Config:  opt.Housekeeping,
			BaseDir: opt.TargetDir,
		})
	default:
		return errors.ErrInvalidAction.Msgf(act)
	}
}

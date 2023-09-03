package tasks

import (
	"context"

	"github.com/vinicius73/thecollector/pkg/cloud"
	"github.com/vinicius73/thecollector/pkg/database"
	"github.com/vinicius73/thecollector/pkg/errors"
)

type Action string

type Options struct {
	Debug       bool
	Database    database.Config
	SyncOptions cloud.SyncOptions
	TargetDir   string
	Datasources []string
}

const (
	ActionDump Action = "dump"
)

func (act Action) Run(ctx context.Context, opt Options) error {
	switch act {
	case ActionDump:
		return Dump(ctx, opt)
	default:
		return errors.ErrInvalidAction.Msgf(act)
	}
}

package tasks

import (
	"context"

	"github.com/vinicius73/thecollector/pkg/housekeeping"
)

func Housekeeping(ctx context.Context, opt housekeeping.Options) error {
	return housekeeping.LocalRemovePreviousDumps(ctx, opt)
}

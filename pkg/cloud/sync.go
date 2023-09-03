package cloud

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/seqsense/s3sync"
)

type SyncOptions struct {
	Bucket      string            `yaml:"bucket"`
	Parallels   int               `yaml:"parallels" default:"2"`
	Credentials RemoteCredentials `yaml:"credentials"`
}

type SyncSource struct {
	BaseDir string
	Dir     string
}

func SyncFiles(ctx context.Context, opt SyncOptions, source SyncSource) error {
	logger := zerolog.Ctx(ctx)

	sess, err := opt.Credentials.DigitalOcean("")
	if err != nil {
		logger.Error().Err(err).Msg("fail to generate cloud credentials")
		return err
	}

	targetDir := opt.Dest(source.Remote())
	sourceDir := source.Dir

	logger.Info().
		Str("source", sourceDir).
		Str("dest", targetDir).
		Msg("syncing with cloud...")

	err = s3sync.New(sess, s3sync.WithParallel(opt.Parallels)).Sync(sourceDir, targetDir)

	if err != nil {
		logger.Error().Err(err).Msg("fail to sync")
		return err
	}

	logger.Info().Msg("sync done")

	return nil
}

func (opt SyncOptions) Dest(dir string) string {
	return fmt.Sprintf("s3://%s", filepath.Join(opt.Bucket, dir))
}

func (ss SyncSource) Remote() string {
	return strings.TrimPrefix(ss.Dir, ss.BaseDir)
}

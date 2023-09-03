package tasks

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	pg "github.com/habx/pg-commands"
	"github.com/rs/zerolog"
	"github.com/vinicius73/thecollector/pkg/checksum"
	"github.com/vinicius73/thecollector/pkg/cloud"
	"github.com/vinicius73/thecollector/pkg/compress"
	"github.com/vinicius73/thecollector/pkg/database"
	"github.com/vinicius73/thecollector/pkg/support"
)

type DumpOptions struct {
	Database  database.Config
	Debug     bool
	TargetDir string
	DBName    string
}

func Dump(ctx context.Context, opt Options) error {
	logger := zerolog.Ctx(ctx)

	if len(opt.Datasources) == 0 {
		logger.Warn().Msg("no datasource to do a backup")
		return nil
	}

	logger.Info().Msgf("running (%v)...", len(opt.Datasources))

	for _, datasource := range opt.Datasources {
		_, err := DumbDB(ctx, DumpOptions{
			Database:  opt.Database,
			TargetDir: opt.TargetDir,
			Debug:     opt.Debug,
			DBName:    datasource,
		})
		if err != nil {
			logger.Error().Err(err).Msg("dump failure")
			continue
		}
	}

	err := cloud.SyncFiles(ctx, opt.SyncOptions, cloud.SyncSource{
		BaseDir: opt.TargetDir,
		Dir:     buildTargetDir(opt.TargetDir),
	})

	return err
}

func DumbDB(ctx context.Context, opt DumpOptions) (string, error) {
	logger := zerolog.Ctx(ctx).With().Str("dbname", opt.DBName).Logger()

	dump, err := pg.NewDump(&pg.Postgres{
		Host:     opt.Database.Host,
		Port:     opt.Database.Port,
		Username: opt.Database.Username,
		Password: opt.Database.Password,
		DB:       opt.Database.Name(opt.DBName),
	})
	if err != nil {
		return "", err
	}

	targetDir := filepath.Join(buildTargetDir(opt.TargetDir), opt.DBName)
	dumpFileName := filepath.Join(targetDir, fmt.Sprintf(`%v.pg_dump`, time.Now().Unix()))

	err = support.EnsureDir(targetDir)

	if err != nil {
		return "", err
	}

	if opt.Debug {
		dump.EnableVerbose()
	}

	dump.SetPath(targetDir + "/")
	dump.SetFileName(filepath.Base(dumpFileName))

	logger.Info().Msgf("generating %s", dumpFileName)

	started := time.Now()

	result := dump.Exec(pg.ExecOptions{
		StreamPrint:       opt.Debug,
		StreamDestination: os.Stdout,
	})
	if result.Error != nil {
		logger.Error().
			Err(result.Error.Err).
			Dur("spended", time.Since(started)).
			Str("command", result.FullCommand).
			Msg(result.Output)

		return dumpFileName, result.Error.Err
	}

	logger.Info().
		Str("file", filepath.Base(dumpFileName)).
		Dur("spended", time.Since(started)).
		Msg("dump done")

	hash, err := checksum.GenerateChecksum(dumpFileName)
	if err != nil {
		logger.Error().Err(err).Msg("fail to generate sum file")
		return dumpFileName, err
	}

	logger.Info().
		Str("file", filepath.Base(dumpFileName)).
		Str("hash", hash).
		Msgf("hash generated")

	compressedFile, err := compress.Compress(ctx, dumpFileName)
	if err != nil {
		logger.Error().Err(err).Msg("fail to compress file")
		return compressedFile, err
	}

	logger.Info().
		Str("file", filepath.Base(compressedFile)).
		Msg("compress done")

	hash, err = checksum.GenerateChecksum(compressedFile)

	if err != nil {
		logger.Error().Err(err).Msg("fail to generate sum file")
		return compressedFile, err
	}

	logger.Info().
		Str("file", filepath.Base(compressedFile)).
		Str("hash", hash).
		Msgf("hash generated")

	if err = os.Remove(dumpFileName); err != nil {
		logger.Warn().
			Err(err).Str("filename", dumpFileName).Msg("fail to remove dump file")
	}

	return compressedFile, nil
}

func buildTargetDir(targetDir string) string {
	return filepath.Join(targetDir, time.Now().Format("2006/01/02"))
}

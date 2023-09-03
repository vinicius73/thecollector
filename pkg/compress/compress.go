package compress

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/klauspost/compress/gzip"
	"github.com/rs/zerolog"
)

func Compress(ctx context.Context, sourceFile string) (string, error) {
	started := time.Now()

	targetName := filepath.Join(filepath.Dir(sourceFile), filepath.Base(sourceFile)+".gz")
	logger := zerolog.Ctx(ctx).With().Str("filename", targetName).Logger()

	logger.Info().Msg("generating gz file...")

	source, err := os.OpenFile(sourceFile, os.O_RDONLY, fs.FileMode(0o644))
	if err != nil {
		return targetName, err
	}

	target, err := os.OpenFile(targetName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fs.FileMode(0o644))
	if err != nil {
		logger.Error().
			Err(err).
			Dur("spended", time.Since(started)).
			Msg("fail create target file")
		return targetName, err
	}

	zw := gzip.NewWriter(target)
	zw.Name = filepath.Base(sourceFile)
	zw.ModTime = time.Now()

	_, err = io.Copy(zw, source)

	if err != nil {
		logger.Error().
			Err(err).
			Dur("spended", time.Since(started)).
			Msg("fail copy source content to gz content")

		return targetName, err
	}

	if err := zw.Close(); err != nil {
		logger.Error().
			Err(err).
			Dur("spended", time.Since(started)).
			Msg("fail close gz file")

		return targetName, err
	}

	logger.Info().
		Dur("spended", time.Since(started)).
		Msg("compact done")

	return targetName, nil
}

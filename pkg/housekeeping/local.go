package housekeeping

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/rs/zerolog"
	"github.com/vinicius73/thecollector/pkg/support"
)

type Options struct {
	Config
	BaseDir string
}

func LocalRemovePreviousDumps(ctx context.Context, opt Options) error {
	logger := zerolog.Ctx(ctx)

	limit := opt.KeepUntil()

	logger.Info().Msgf("removing dumps older than %v", limit)

	dumps, err := listAllDumps(opt.BaseDir)
	if err != nil {
		return err
	}

	ch := make(chan DumpDir, 2)

	var wg sync.WaitGroup

	for i := 0; i < opt.Workers; i++ {
		wg.Add(1)
		go func() {
			for dump := range ch {
				if dump.Date.Before(limit) {
					rel, _ := filepath.Rel(opt.BaseDir, dump.Path)
					logger.Warn().Msgf("removing %v", rel)

					err := os.RemoveAll(dump.Path)
					if err != nil {
						logger.Error().Err(err).Msgf("failed to remove %v", rel)
					}

				}
			}
			wg.Done()
		}()
	}

	for _, dump := range dumps {
		if dump.Date.Before(limit) {
			ch <- dump
		}
	}

	close(ch)

	wg.Wait()

	err = removeEmptyFolders(opt.BaseDir)

	if err != nil {
		return err
	}

	logger.Info().Msg("dump removal completed")

	return nil
}

func removeEmptyFolders(baseDir string) error {
	dirs, err := support.ListDirs(baseDir, 0, 3)
	if err != nil {
		return err
	}

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i] > dirs[j]
	})

	for _, dir := range dirs {
		isEmpty, err := support.IsEmptyDir(dir)
		if err != nil {
			return err
		}

		if isEmpty {
			err := os.RemoveAll(dir)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

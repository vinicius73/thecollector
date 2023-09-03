package config

import (
	"os"
	"path/filepath"

	"github.com/creasty/defaults"
	"github.com/vinicius73/thecollector/pkg/errors"
	"github.com/vinicius73/thecollector/pkg/support"
	"gopkg.in/yaml.v3"
)

func Load(fileName string) (*App, error) {
	cfg := App{}

	var err error

	if fileName == "" {
		if fileName, err = tryToFindConfigFile(); err != nil {
			return nil, err
		}
	}

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, errors.ErrFailToLoadConfig.Wrap(err)
	}

	// expand env vars
	bytes = []byte(os.ExpandEnv(string(bytes)))

	if err = yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, errors.ErrFailToLoadConfig.Wrap(err)
	}

	return applyDefaults(cfg)
}

func applyDefaults(cfg App) (*App, error) {
	err := defaults.Set(&cfg)
	if err != nil {
		return &cfg, err
	}

	if cfg.TargetDir == "" {
		cfg.TargetDir = support.GetEnv("THECOLLECTOR_TARGET_DIR", "./dumps")
	}

	if !filepath.IsAbs(cfg.TargetDir) {
		if pwd, err := os.Getwd(); err == nil {
			cfg.TargetDir = filepath.Join(pwd, cfg.TargetDir)
		} else {
			return &cfg, err
		}
	}

	return &cfg, nil
}

func tryToFindConfigFile() (string, error) {
	env := support.GetEnv("THECOLLECTOR_CONFIG_FILE", "")
	pwd, err := os.Getwd()
	if err != nil {
		return "", errors.ErrFailToLoadConfig.Wrap(err)
	}

	possiblePaths := []string{
		filepath.Join(pwd, "thecollector.yml"),
		"./thecollector.yml",
		filepath.Join(support.GetBinDirPath(), "thecollector.yml"),
	}

	if env != "" {
		// put env in the first position
		possiblePaths = append([]string{env}, possiblePaths...)
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", errors.ErrConfigNotFound
}

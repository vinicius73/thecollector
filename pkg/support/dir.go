package support

import (
	"os"
	"path/filepath"
)

func GetBinDirPath() string {
	execPath, _ := os.Executable()
	return filepath.Dir(execPath)
}

func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0o700)

		if err != nil {
			return err
		}
	}

	return nil
}

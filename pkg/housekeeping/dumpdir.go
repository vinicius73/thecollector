package housekeeping

import (
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/vinicius73/thecollector/pkg/support"
)

type DumpDir struct {
	Path     string
	Database string
	Date     time.Time
}

func parseDumpDir(path string) DumpDir {
	parts := support.Reverse(strings.Split(path, "/"))

	database := parts[0]

	year, _ := strconv.Atoi(parts[3])
	month, _ := strconv.Atoi(parts[2])
	day, _ := strconv.Atoi(parts[1])

	return DumpDir{
		Path:     path,
		Database: database,
		Date:     time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local),
	}
}

func listAllDumps(baseDir string) ([]DumpDir, error) {
	baseDir = path.Clean(baseDir)
	dirs, err := support.ListDirs(baseDir, 0, 4)
	var folders []DumpDir

	if err != nil {
		return nil, err
	}

	// filter dump dirs
	for _, dir := range dirs {
		if IsDumpDir(strings.TrimPrefix(dir, baseDir)) {
			folders = append(folders, parseDumpDir(dir))
		}
	}

	return folders, nil
}

func IsDumpDir(dir string) bool {
	// remove leading slash and split
	parts := strings.Split(strings.TrimPrefix(dir, "/"), "/")

	// YYYY/MM/DD/database
	if len(parts) != 4 {
		return false
	}

	// check if the first part is a year
	if _, err := strconv.Atoi(parts[0]); err != nil {
		return false
	}

	// check if the second part is a month
	if _, err := strconv.Atoi(parts[1]); err != nil {
		return false
	}

	// check if the third part is a day
	if _, err := strconv.Atoi(parts[2]); err != nil {
		return false
	}

	// check if the fourth part is a database
	if len(parts[3]) == 0 {
		return false
	}

	return true
}

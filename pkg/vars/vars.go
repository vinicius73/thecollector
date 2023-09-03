package vars

import (
	"fmt"
	"os"
)

var (
	version   string
	buildDate string
	commit    string
)

func Version() string {
	return version
}

func VersionVerbose() string {
	return fmt.Sprintf("Version %s\nRevision %s\nBuild at %s", Version(), Commit(), BuildDate())
}

func BuildDate() string {
	return buildDate
}

func Commit() string {
	return commit
}

func Host() string {
	h, _ := os.Hostname()

	if h == "" {
		h = "unknown"
	}

	return h
}

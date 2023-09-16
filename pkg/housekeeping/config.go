package housekeeping

import "time"

type Config struct {
	KeepLocal time.Duration `yaml:"keep_local" default:"336h"`
	Workers   int           `yaml:"workers" default:"2"`
}

func (c Config) KeepUntil() time.Time {
	date := time.Now().Add(-c.KeepLocal)
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

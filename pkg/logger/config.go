package logger

type Config struct {
	Level  string `yaml:"level" default:"debug"`
	Format string `yaml:"format" default:"text"`
}

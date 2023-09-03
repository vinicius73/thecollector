package database

type Config struct {
	Host         string `yaml:"host" default:"localhost"`
	Port         int    `yaml:"port" default:"5432"`
	Username     string `yaml:"username" default:"theconnector"`
	Password     string `yaml:"password"`
	DBNamePrefix string `yaml:"dbname_prefix"`
}

func (c Config) Name(name string) string {
	if c.DBNamePrefix != "" {
		return c.DBNamePrefix + name
	}

	return name
}

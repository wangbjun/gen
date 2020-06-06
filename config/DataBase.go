package config

var DBConfig map[string]map[string]string

func init() {
	DBConfig = map[string]map[string]string{
		"default": {
			"dialect":      Conf.Section("DB").Key("Dialect").String(),
			"dsn":          Conf.Section("DB").Key("DSN").String(),
			"maxIdleConns": Conf.Section("DB").Key("MAX_IDLE_CONN").String(),
			"maxOpenConns": Conf.Section("DB").Key("MAX_OPEN_CONN").String(),
		},
		"user": {
			"dialect":      Conf.Section("DB").Key("Dialect").String(),
			"dsn":          Conf.Section("DB").Key("DSN").String(),
			"maxIdleConns": Conf.Section("DB").Key("MAX_IDLE_CONN").String(),
			"maxOpenConns": Conf.Section("DB").Key("MAX_OPEN_CONN").String(),
		},
	}
}

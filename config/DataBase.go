package config

var DBConfig map[string]map[string]string

func init() {
	DBConfig = map[string]map[string]string{
		"default": {
			"dialect":      Conf.String("DB_Dialect"),
			"host":         Conf.String("DB_HOST"),
			"port":         Conf.String("DB_PORT"),
			"database":     Conf.String("DB_DATABASE"),
			"username":     Conf.String("DB_USERNAME"),
			"password":     Conf.String("DB_PASSWORD"),
			"charset":      Conf.String("DB_CHARSET"),
			"maxIdleConns": Conf.String("DB_MAX_IDLE_CONN"),
			"maxOpenConns": Conf.String("DB_MAX_OPEN_CONN"),
		},
	}
}

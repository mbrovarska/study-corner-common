package config

type AppConfig struct {
	ServiceName string
	HTTPPort    int
	DB_DSN      string
	ENV         string
	LogLevel    string
}

type Loader interface {
	Load() (AppConfig, error)
}
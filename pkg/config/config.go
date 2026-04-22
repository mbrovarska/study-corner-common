package config

type AppConfig struct {
	ServiceName string
	HTTPPort    int
	DBDSN       string
	Env         string
	LogLevel    string
}


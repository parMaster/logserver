package logserver

// Config ...
type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	MqUser      string `toml:"mq_user"`
	MqPassword  string `toml:"mq_password"`
	MqClientId  string `toml:"mq_client_id"`
	MqBrokerURL string `toml:"mq_broker_url"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}

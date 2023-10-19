package config

import (
	"github.com/BurntSushi/toml"
)

// Config ...
type Config struct {
	BindAddr     string `toml:"bind_addr"`
	LogLevel     string `toml:"log_level"`
	DatabaseKind string `toml:"database_kind"`
	DatabaseURL  string `toml:"database_url"`
	MqUser       string `toml:"mq_user"`
	MqPassword   string `toml:"mq_password"`
	MqClientId   string `toml:"mq_client_id"`
	MqBrokerURL  string `toml:"mq_broker_url"`
	MqRootTopic  string `toml:"mq_root_topic"`
	CollectRaw   bool   `toml:"collect_raw"`
}

// NewConfig ...
func NewConfig(filepath string) (*Config, error) {
	c := Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}

	_, err := toml.DecodeFile(filepath, &c)

	return &c, err
}

// ToDo rewrite to yaml to load modules and configure them

package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// server:
//
//	bind_addr: ":8088"
//	dbg: true # enable debug mode, can be set to true with --dbg flag when running the service
type Server struct {
	BindAddr string `yaml:"bind_addr"`
	Dbg      bool   `yaml:"dbg"`
}

// # mqtt credentials
// mqtt:
//
//	mq_user: foo
//	mq_password: bar
//	mq_client_id: baz
//	mq_broker_url: ssl://mqtt.foobar:8883
//	mq_root_topic: "#"
type Mqtt struct {
	MqUser      string `yaml:"mq_user"`
	MqPassword  string `yaml:"mq_password"`
	MqClientId  string `yaml:"mq_client_id"`
	MqBrokerURL string `yaml:"mq_broker_url"`
	MqRootTopic string `yaml:"mq_root_topic"`
}

// storage:
//
//	type: sqlite
//	database_url: file:./mqttdata.db?mode=rwc
type Storage struct {
	Type string `yaml:"type"` // Type of storage to use. Currently supported: sqlite
	Path string `yaml:"path"` // Path to the database file
}

// collect_raw: false
type Config struct {
	Server     Server  `yaml:"server"`
	Mqtt       Mqtt    `yaml:"mqtt"`
	Storage    Storage `yaml:"storage"`
	CollectRaw bool    `yaml:"collect_raw"`
}

// NewConfig creates a new Config from the given file
func NewConfig(fname string) (*Config, error) {
	c := &Config{}
	data, err := os.ReadFile(fname)
	if err != nil {
		log.Printf("[ERROR] can't read config %s: %e", fname, err)
		return nil, fmt.Errorf("can't read config %s: %w", fname, err)
	}
	if err = yaml.Unmarshal(data, &c); err != nil {
		log.Printf("[ERROR] failed to parse config %s: %e", fname, err)
		return nil, fmt.Errorf("failed to parse config %s: %w", fname, err)
	}
	// log.Printf("[DEBUG] config: %+v", p)
	return c, nil
}

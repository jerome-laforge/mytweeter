package config

import (
	"sync"

	"github.com/go-ini/ini"
)

var (
	conf *Config
	once sync.Once
	err  error
)

type Config struct {
	Cassandra struct {
		Cluster  []string
		Keyspace string
	}
	Web struct {
		Address string
	}
	Hystrix struct {
		Timeout                string
		MaxConcurrentRequests  int
		RequestVolumeThreshold int
		SleepWindow            int
		ErrorPercentThreshold  int
	}
	Endless struct {
		DefaultHammerTime string
	}
}

func GetConfig() (*Config, error) {
	once.Do(func() {
		conf = new(Config)
		err = ini.MapTo(conf, "conf.ini")
	})
	return conf, err
}

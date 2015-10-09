package dao

import (
	"config"
	"sync"

	"github.com/gocql/gocql"
)

var (
	instance *gocql.ClusterConfig
	once     sync.Once
	err      error
)

func NewSession() (*gocql.Session, error) {
	once.Do(func() {
		var conf *config.Config
		conf, err = config.GetConfig()
		if err != nil {
			return
		}
		instance = gocql.NewCluster(conf.Cassandra.Cluster...)
		instance.Keyspace = conf.Cassandra.Keyspace
		instance.Consistency = gocql.Quorum
	})

	if err != nil {
		return nil, err
	}

	return instance.CreateSession()
}

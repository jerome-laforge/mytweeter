package dao

import (
	"config"
	"sync"

	"github.com/gocql/gocql"
)

var (
	session *gocql.Session
	once    sync.Once
	err     error
)

func GetSession() (*gocql.Session, error) {
	once.Do(func() {
		var conf *config.Config
		conf, err = config.GetConfig()
		if err != nil {
			return
		}
		cluster := gocql.NewCluster(conf.Cassandra.Cluster...)
		cluster.Keyspace = conf.Cassandra.Keyspace
		cluster.Consistency = gocql.Quorum
		session, err = cluster.CreateSession()
	})
	return session, err
}

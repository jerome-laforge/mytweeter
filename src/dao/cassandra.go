package dao

import (
	"sync"

	"github.com/gocql/gocql"
)

var (
	instance *gocql.ClusterConfig
	once     sync.Once
)

func NewSession() (*gocql.Session, error) {
	once.Do(func() {
		instance = gocql.NewCluster("127.0.0.1")
		instance.Keyspace = "example"
		instance.Consistency = gocql.Quorum
	})

	return instance.CreateSession()
}

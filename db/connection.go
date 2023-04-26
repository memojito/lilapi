package db

import (
	"github.com/gocql/gocql"
	"time"
)

func NewSession() (*gocql.Session, error) {
	// Define the Cassandra cluster configuration
	clusterConfig := gocql.NewCluster("127.0.0.1:9042")
	clusterConfig.Keyspace = "lilapi"
	clusterConfig.Consistency = gocql.Quorum
	clusterConfig.ProtoVersion = 4
	clusterConfig.ConnectTimeout = time.Second * 10

	// Create the session object
	session, err := clusterConfig.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

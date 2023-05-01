package db

import (
	"github.com/gocql/gocql"
	"log"
	"time"
)

func NewSession() (*gocql.Session, error) {
	// Define the Cassandra cluster configuration
	clusterConfig := gocql.NewCluster("cassandra")
	clusterConfig.Keyspace = "lilapi"
	clusterConfig.Consistency = gocql.Quorum
	clusterConfig.ProtoVersion = 4
	clusterConfig.ConnectTimeout = time.Second * 10

	var createdSession *gocql.Session
	for {
		session, err := clusterConfig.CreateSession()
		if err == nil {
			createdSession = session
			break
		}
		log.Printf("CreateSession: %v", err)
		time.Sleep(time.Second)
	}
	log.Printf("Connected OK")

	return createdSession, nil
}

package database

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var Driver neo4j.DriverWithContext

func Init(uri, username, password string) error {
	var err error
	Driver, err = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	return err
}

func Close(ctx context.Context) {
	if Driver != nil {
		err := Driver.Close(ctx)
		if err != nil {
			log.Print("Failed to close Neo4j driver: ", err)
		}
	}
}

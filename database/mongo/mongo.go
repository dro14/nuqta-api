package mongo

import (
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Mongo struct {
	users *mongo.Collection
	posts *mongo.Collection
}

func New() *Mongo {
	uri, ok := os.LookupEnv("MONGO_URI")
	if !ok {
		log.Fatal("mongodb uri is not specified")
	}

	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal("can't connect to mongodb: ", err)
	}

	return &Mongo{
		users: client.Database("nuqta").Collection("users"),
		posts: client.Database("nuqta").Collection("posts"),
	}
}

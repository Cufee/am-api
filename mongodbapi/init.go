package mongodbapi

import (
	"context"
	"log"
	"time"

	"github.com/cufee/am-api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Collections
var userDataCollection *mongo.Collection
var intentsCollection *mongo.Collection
var ctx = context.TODO()

func init() {
	// Conenct to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Println("Panic in mongoapi/init")
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Panic in mongoapi/init")
		panic(err)
	}
	log.Println("Successfully connected and pinged.")

	// Collections
	userDataCollection = client.Database("webapp").Collection("users")
	intentsCollection = client.Database("webapp").Collection("intents")
}

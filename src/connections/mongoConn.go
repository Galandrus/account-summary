package connections

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	dbDefaultTimeout        = 30 * time.Second
	dbDefaultConnectTimeout = 10 * time.Second
	dbDefaultSocketTimeout  = 10 * time.Minute
	dbDefaultPoolSize       = uint64(10)
)

func NewMongoConn(uri string) *mongo.Client {
	if uri == "" {
		panic("MONGO_URI is not set")
	}
	clientOptions := options.ClientOptions{
		ConnectTimeout:  &dbDefaultConnectTimeout,
		SocketTimeout:   &dbDefaultSocketTimeout,
		MaxConnIdleTime: &dbDefaultTimeout,
		MaxPoolSize:     &dbDefaultPoolSize,
		MinPoolSize:     &dbDefaultPoolSize,
	}

	ctxTimeout, cancel := context.WithTimeout(context.TODO(), dbDefaultConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctxTimeout, clientOptions.ApplyURI(uri))

	if err != nil {
		panic(fmt.Sprintf("mongoDB error in client config: %s", err.Error()))
	}

	if err = client.Ping(ctxTimeout, readpref.Primary()); err != nil {
		panic(fmt.Sprintf("mongoDB error in client connection: %s", err.Error()))
	}
	return client
}

func GetCollection(client *mongo.Client, database, collection string) *mongo.Collection {
	return client.Database(database).Collection(collection)
}

func Disconnect(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		return
	}
}

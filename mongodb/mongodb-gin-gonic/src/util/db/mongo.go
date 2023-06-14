package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoDbUtil struct {
	srv            string
	dbName         string
	collectionName string
	ctx            context.Context
}

func NewMongoDbUtil(srv, dbName, collectionName string) *MongoDbUtil {
	return &MongoDbUtil{
		srv:            srv,
		dbName:         dbName,
		collectionName: collectionName,
		ctx:            context.Background(),
	}

}

func (o *MongoDbUtil) GetCollection() (client *mongo.Client, col *mongo.Collection) {
	client, err := o.connect()
	if err != nil {
		log.Println(err)
		err = errors.New("fail connect to data")
		return
	}
	col = client.Database(o.dbName).Collection(o.collectionName)
	return

}

func (o *MongoDbUtil) connect() (client *mongo.Client, err error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI(o.srv)
	clientOptions.SetMaxPoolSize(200)

	client, err = mongo.Connect(o.ctx, clientOptions)
	if err != nil {
		log.Println(err)
		err = errors.New("fail connect to data")
		return
	}
	return
}

package db

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDb *mongo.Database

func GetMongoDB() (*mongo.Database, error) {
	if mongoDb == nil {
		_db, err := getMongoDbConnection()
		if err != nil {
			return nil, err
		}
		mongoDb = _db
	}

	return mongoDb, nil
}

func getMongoDbConnection() (*mongo.Database, error) {
	mongoUrl := os.Getenv("MONGO_URL")
	mongoDatabaseName := os.Getenv("MONGO_DATABASE_NAME")
	if mongoUrl == "" || mongoDatabaseName == "" {
		return nil, errors.New("MONGO_URL or MONGO_DATABASE_NAME is not set")
	}

	loggerOptions := options.Logger().SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
	clientOptions := options.Client().ApplyURI(mongoUrl).SetLoggerOptions(loggerOptions)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	return client.Database(mongoDatabaseName), nil
}

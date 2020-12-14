package database

import "go.mongodb.org/mongo-driver/mongo"

type CollectionsStruct struct {
	Image    string
	Password string
	User     string
}

type MongoInstance struct {
	Client   *mongo.Client
	Database *mongo.Database
}

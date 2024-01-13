package models

import "go.mongodb.org/mongo-driver/mongo"

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    int32   `json:"age"`
}

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var MG MongoInstance

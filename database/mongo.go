package database

import (
	"context"
	"github.com/komron-dev/hrms-fiber/config"
	"github.com/komron-dev/hrms-fiber/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	defer cancel()
	if err != nil {
		return err
	}
	db := client.Database(config.DBName)

	models.MG = models.MongoInstance{
		Client: client,
		Db:     db,
	}

	return nil
}

package db

import (
	"6_gopro/config"

	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database

func Connectdb() {

	if config.Conf.Db.Username == "" || config.Conf.Db.Password == "" {

		url := fmt.Sprintf("mongodb://%s:%s", config.Conf.Db.Host, config.Conf.Db.Port)

		client, err := mongo.NewClient(options.Client().ApplyURI(url))
		if err != nil {
			panic(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = client.Connect(ctx)
		if err != nil {
			panic(err)
		}

		Db = client.Database(config.Conf.Db.Name)

		_, err = Db.ListCollectionNames(context.Background(), bson.M{})
		if err != nil {
			panic(err)
		}

	} else {

		credential := options.Credential{
			Username: config.Conf.Db.Username,
			Password: config.Conf.Db.Password,
		}

		url := fmt.Sprintf("mongodb://%s:%s", config.Conf.Db.Host, config.Conf.Db.Port)

		client, err := mongo.NewClient(options.Client().ApplyURI(url).SetAuth(credential))
		if err != nil {
			panic(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = client.Connect(ctx)
		if err != nil {
			panic(err)
		}

		Db = client.Database(config.Conf.Db.Name)

		_, err = Db.ListCollectionNames(context.Background(), bson.M{})
		if err != nil {
			panic(err)
		}

	}

}

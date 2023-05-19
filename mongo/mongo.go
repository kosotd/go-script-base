package mongo

import (
	"context"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	db *mongo.Client
}

func Init() *mongo.Client {
	opts := options.Client()
	opts = opts.ApplyURI(viper.GetString("mongo_uri"))
	opts = opts.SetMaxPoolSize(viper.GetUint64("mongo_max_pool_size"))

	client, err := mongo.NewClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	return client
}

func (r *Mongo) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	defer r.db.Disconnect(ctx)
}

func (r *Mongo) Instance() *mongo.Client {
	return r.db
}

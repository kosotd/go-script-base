package redis

import (
	"context"
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"crypto/tls"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	RDB *redis.Client
}

var RDB *redis.Client

func init() {
	pflag.String("redis_host", "127.0.0.1:6379", "Redis server address")
	pflag.String("redis_username", "", "Redis username")
	pflag.String("redis_password", "", "Redis password")
	pflag.Int("redis_db", 0, "Redis DB id")
	pflag.Bool("redis_use_tls", false, "Use TLS or not")

	ctx := context.Background()

	var tlsConfig *tls.Config
	if viper.GetBool("redis_use_tls") {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:      viper.GetString("redis_host"),
		Username:  viper.GetString("redis_username"),
		Password:  viper.GetString("redis_password"),
		DB:        viper.GetInt("redis_db"),
		TLSConfig: tlsConfig,
	})

	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}
}

func Close() {
	RDB.Close()
}

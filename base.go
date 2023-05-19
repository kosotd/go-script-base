package base

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.String("log_level", "INFO", "log level")

	pflag.String("postgres_uri", "postgres://root:1234@localhost:5432/nft_pipeline", "Postgres connection string")
	pflag.Int("db_max_open_conns", 0, "db driver max open connections (0 - that is no limit on the number)")
	pflag.Int("db_max_idle_conns", 2, "db driver max idle connections")
	pflag.Duration("db_conn_max_lifetime", 5*time.Minute, "db driver connection max lifetime")
	pflag.String("application_name", "application", "application name")

	pflag.String("redis_host", "127.0.0.1:6379", "Redis server address")
	pflag.String("redis_username", "", "Redis username")
	pflag.String("redis_password", "", "Redis password")
	pflag.Int("redis_db", 0, "Redis DB id")
	pflag.Bool("redis_use_tls", false, "Use TLS or not")

	pflag.Duration("sleep_interval", 0, "sleep interval")

	pflag.String("mongo_uri", "mongo://localhost:27017", "MongoDB url")
	pflag.Uint64("mongo_max_pool_size", 300, "MongoDB max pool size")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
}

type Handler interface {
	Handle(ctx context.Context) error
}

func Run(handler Handler) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if viper.GetDuration("sleep_interval") <= 0 {
		if err := cycle(ctx, handler); err != nil {
			return err
		}
	} else {
		ticker := time.NewTicker(viper.GetDuration("sleep_interval"))
		defer ticker.Stop()

		if err := handler.Handle(ctx); err != nil {
			return err
		}

		if err := cycleTicker(ctx, ticker, handler); err != nil {
			return err
		}
	}

	return nil
}

func cycle(ctx context.Context, handler Handler) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := handler.Handle(ctx); err != nil {
				return err
			}
		}
	}
}

func cycleTicker(ctx context.Context, ticker *time.Ticker, handler Handler) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := handler.Handle(ctx); err != nil {
				return err
			}
		}
	}
}

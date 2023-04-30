package postgres

import (
	"context"
	"log"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	DB   *sqlx.DB
	Pool *pgxpool.Pool
}

var DB *sqlx.DB
var Pool *pgxpool.Pool

func init() {
	pflag.String("postgres_uri", "postgres://root:1234@localhost:5432/nft_pipeline", "Postgres connection string")
	pflag.Int("db_max_open_conns", 0, "db driver max open connections (0 - that is no limit on the number)")
	pflag.Int("db_max_idle_conns", 2, "db driver max idle connections")
	pflag.Duration("db_conn_max_lifetime", 5*time.Minute, "db driver connection max lifetime")
	pflag.String("application_name", "application", "application name")

	ctx := context.Background()
	postgresUri := viper.GetString("postgres_uri")
	var err error

	DB, err = sqlx.ConnectContext(ctx, "pgx", postgresUri)
	if err != nil {
		log.Fatal(err)
	}

	Pool, err = pgxpool.New(ctx, postgresUri)
	if err != nil {
		log.Fatal(err)
	}

	DB.SetMaxOpenConns(viper.GetInt("db_max_open_conns"))
	DB.SetMaxIdleConns(viper.GetInt("db_max_idle_conns"))
	DB.SetConnMaxLifetime(viper.GetDuration("db_conn_max_lifetime"))

	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}

	if _, err := DB.Exec(fmt.Sprintf("SET application_name = '%s'", viper.GetString("application_name"))); err != nil {
		log.Fatal(err)
	}
}

func Close() {
	DB.Close()
	Pool.Close()
}

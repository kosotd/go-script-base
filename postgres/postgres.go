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
	db   *sqlx.DB
	pool *pgxpool.Pool
}

func Init() *Postgres {
	pflag.String("postgres_uri", "postgres://root:1234@localhost:5432/nft_pipeline", "Postgres connection string")
	pflag.Int("db_max_open_conns", 0, "db driver max open connections (0 - that is no limit on the number)")
	pflag.Int("db_max_idle_conns", 2, "db driver max idle connections")
	pflag.Duration("db_conn_max_lifetime", 5*time.Minute, "db driver connection max lifetime")
	pflag.String("application_name", "application", "application name")

	ctx := context.Background()
	postgresUri := viper.GetString("postgres_uri")

	db, err := sqlx.ConnectContext(ctx, "pgx", postgresUri)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, postgresUri)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(viper.GetInt("db_max_open_conns"))
	db.SetMaxIdleConns(viper.GetInt("db_max_idle_conns"))
	db.SetConnMaxLifetime(viper.GetDuration("db_conn_max_lifetime"))

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(fmt.Sprintf("SET application_name = '%s'", viper.GetString("application_name"))); err != nil {
		log.Fatal(err)
	}

	return &Postgres{db: db, pool: pool}
}

func (p *Postgres) Close() {
	defer p.db.Close()
	defer p.pool.Close()
}

func (p *Postgres) Instance() (*sqlx.DB, *pgxpool.Pool) {
	return p.db, p.pool
}

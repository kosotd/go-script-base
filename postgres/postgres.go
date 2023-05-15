package postgres

import (
	"context"
	"log"

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

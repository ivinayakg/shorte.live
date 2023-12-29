package models

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func ConnectToDB() (*bun.DB, error) {
	maxOpenConnectionsStr := os.Getenv("DB_MAX_OPEN_CONNECTIONS")
	maxOpenConnections, err := strconv.Atoi(maxOpenConnectionsStr)

	if err != nil {
		return nil, err
	}

	// dsn := "postgres://postgres:@localhost:5432/test?sslmode=disable"
	// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
	dsn := os.Getenv("DB_URL")
	pgDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	pgDb.SetMaxOpenConns(maxOpenConnections)

	db := bun.NewDB(pgDb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db, nil
}

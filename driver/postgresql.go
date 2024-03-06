package driver

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/assyatier21/learn-distributed-trx/config"
	_ "github.com/jackc/pgx"
)

func NewPostgreSQL(cfg config.DBConfig) (*sql.DB, error) {
	dsn := cfg.GetDSN()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return db, err
	}

	return db, nil
}

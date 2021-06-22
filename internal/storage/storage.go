package storage

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool
var err error

// urlExample := "postgres://username:password@localhost:5432/database_name"
var Url  = ""
func init() {
	pool, err = pgxpool.Connect(context.Background(), Url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to database:", err)
		os.Exit(1)
	}
}

func GetCurrentPeriod(period_match string) (int, error) {
	var period int
	err = pool.QueryRow(context.Background(),
		"select id from rsj.period where id = $1 and is_open=True limit 1", period_match).Scan(&period)
	if err != nil {
		log.Println(err)
		return period, err
	}
	return period, err
}


package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func NewConn(url string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Panic(err)
	}
	return conn, nil
}

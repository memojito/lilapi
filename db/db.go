package db

import (
	"github.com/gocql/gocql"
	"log"
)

type Session struct {
	Session *gocql.Session
}

type Transaction struct {
	ID    gocql.UUID
	Name  string
	Value int
}

func (session Session) SaveTransaction(name string, price int) {
	id, _ := gocql.RandomUUID()
	if err := session.Session.Query(`INSERT INTO transaction (id, name, value) VALUES (?, ?, ?)`,
		id, name, price).Exec(); err != nil {
		log.Println(err)
		return
	}
}

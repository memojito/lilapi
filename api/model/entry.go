package model

import "github.com/gocql/gocql"

type Transaction struct {
	ID    gocql.UUID
	Name  string
	Value int
}

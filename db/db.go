package db

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

type Session struct {
	Session *gocql.Session
}

type Conn struct {
	Conn *pgx.Conn
}

type Transaction struct {
	ID           int64
	Name         string
	Value        int
	UserID       int64     `db:"user_id"`
	CreationDate time.Time `db:"creation_date"`
}

type User struct {
	ID        int64
	FirstName string
	LastName  string
}

func (conn Conn) SaveTransaction(name string, price int, userID int64) {
	if _, err := conn.Conn.Exec(context.Background(), `INSERT INTO transaction (name, value, user_id, creation_date) VALUES ($1, $2, $3, $4)`,
		name, price, userID, time.Now()); err != nil {
		log.Println(err)
		return
	}
}

func (conn Conn) GetTransaction(id int64) (Transaction, error) {
	transaction := Transaction{}
	if err := conn.Conn.QueryRow(context.Background(), `SELECT id, name, value, user_id FROM transaction WHERE id = $1`, id).Scan(&transaction.ID,
		&transaction.Name, &transaction.Value, &transaction.Value); err != nil {
		log.Println(err)
		return Transaction{}, err
	}
	return transaction, nil
}

func (conn Conn) GetAllTransactionsByUserIdAndDate(userID int64, date time.Time) ([]Transaction, error) {
	q := `SELECT id, name, value, user_id, creation_date FROM transaction WHERE (user_id = $1 AND creation_date = $2)`

	rows, err := conn.Conn.Query(context.Background(), q, userID, date)
	if err != nil {
		log.Printf("Failed query: %s\n", err)
		return nil, err
	}

	transactions, err := pgx.CollectRows(rows, pgx.RowToStructByName[Transaction])
	if err != nil {
		log.Printf("Failed collecting rows: %s\n", err)
		return nil, err
	}

	for _, t := range transactions {
		log.Println(t)
	}

	return transactions, nil
}

func (conn Conn) SaveUser(id int64, firstName string, lastName string) {
	if _, err := conn.Conn.Exec(context.Background(), `INSERT INTO teleuser (id, first_name, last_name) VALUES ($1, $2, $3)`,
		id, firstName, lastName); err != nil {
		log.Println(err)
		return
	}
}

func (conn Conn) GetUser(id int64) (User, error) {
	user := User{}
	if err := conn.Conn.QueryRow(context.Background(), `SELECT id, first_name, last_name FROM teleuser WHERE id = $1`, id).Scan(&user.ID,
		&user.FirstName, &user.LastName); err != nil {
		log.Println(err)
		return User{}, err
	}
	return user, nil
}

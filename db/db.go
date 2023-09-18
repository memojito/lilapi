package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type Conn struct {
	Conn *pgx.Conn
}

type Transaction struct {
	ID           int64
	Name         string
	Value        int
	UserID       int64     `db:"user_id"`
	CreationDate time.Time `db:"creation_date"`
	CategoryID   int64     `db:"category_id"`
}

type Category struct {
	ID     int64
	Name   string
	UserID int64 `db:"user_id"`
}

type User struct {
	ID        int64
	FirstName string
	LastName  string
}

func (conn *Conn) SaveTransaction(name string, price int, userID int64, categoryID int64) {
	q := "INSERT INTO transaction (name, value, user_id, creation_date, category_id) VALUES ($1, $2, $3, $4, $5)"

	if _, err := conn.Conn.Exec(context.Background(), q, name, price, userID, time.Now(), categoryID); err != nil {
		log.Println(err)
		return
	}
}

func (conn *Conn) GetTransaction(id int64) (Transaction, error) {
	transaction := Transaction{}
	q := "SELECT id, name, value, user_id, creation_date, category_id FROM transaction WHERE id = $1"

	if err := conn.Conn.QueryRow(context.Background(), q, id).Scan(&transaction.ID,
		&transaction.Name, &transaction.Value, &transaction.Value); err != nil {
		log.Println(err)
		return Transaction{}, err
	}
	return transaction, nil
}

func (conn *Conn) GetTransactions(userID int64, date time.Time) ([]Transaction, error) {
	q := "SELECT id, name, value, user_id, creation_date, category_id FROM transaction WHERE (user_id = $1 AND creation_date = $2)"

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

func (conn *Conn) SaveCategory(name string, userID int64) {
	q := "INSERT INTO category (name, user_id) VALUES ($1, $2)"

	if _, err := conn.Conn.Exec(context.Background(), q, name, userID); err != nil {
		log.Println(err)
		return
	}
}

func (conn *Conn) GetCategory(id int64) (Category, error) {
	category := Category{}
	q := "SELECT id, name, user_id FROM category WHERE id = $1"

	if err := conn.Conn.QueryRow(context.Background(), q, id).Scan(&category.ID,
		&category.Name, &category.UserID); err != nil {
		log.Println(err)
		return Category{}, err
	}
	return category, nil
}

func (conn *Conn) GetCategoryByName(name string, userID int64) (Category, error) {
	category := Category{}
	q := "SELECT id, name, user_id FROM category WHERE (name = $1 AND user_id = $2)"

	if err := conn.Conn.QueryRow(context.Background(), q, name, userID).Scan(&category.ID,
		&category.Name, &category.UserID); err != nil {
		log.Println(err)
		return Category{}, err
	}
	return category, nil
}

func (conn *Conn) GetCategories(userID int64) ([]Category, error) {
	q := "SELECT id, name, user_id FROM category WHERE user_id = $1"

	rows, err := conn.Conn.Query(context.Background(), q, userID)
	if err != nil {
		log.Printf("Failed query: %s\n", err)
		return nil, err
	}

	categories, err := pgx.CollectRows(rows, pgx.RowToStructByName[Category])
	if err != nil {
		log.Printf("Failed collecting rows: %s\n", err)
		return nil, err
	}

	return categories, nil
}

func (conn *Conn) SaveUser(id int64, firstName string, lastName string) {
	q := "INSERT INTO teleuser (id, first_name, last_name) VALUES ($1, $2, $3)"

	if _, err := conn.Conn.Exec(context.Background(), q, id, firstName, lastName); err != nil {
		log.Println(err)
		return
	}
}

func (conn *Conn) GetUser(id int64) (User, error) {
	user := User{}
	q := "SELECT id, first_name, last_name FROM teleuser WHERE id = $1"

	if err := conn.Conn.QueryRow(context.Background(), q, id).Scan(&user.ID,
		&user.FirstName, &user.LastName); err != nil {
		log.Println(err)
		return User{}, err
	}
	return user, nil
}

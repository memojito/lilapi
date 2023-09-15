package db

import (
	"github.com/gocql/gocql"
	"log"
)

type Session struct {
	Session *gocql.Session
}

type Transaction struct {
	ID     gocql.UUID
	Name   string
	Value  int
	UserID int64
}

type User struct {
	ID        int64
	FirstName string
	LastName  string
}

func (session Session) SaveTransaction(name string, price int, userID int64) {
	id, _ := gocql.RandomUUID()
	if err := session.Session.Query(`INSERT INTO transaction (id, name, value, user_id) VALUES (?, ?, ?, ?)`,
		id, name, price, userID).Exec(); err != nil {
		log.Println(err)
		return
	}
}

func (session Session) GetTransaction(id gocql.UUID) (Transaction, error) {
	transaction := Transaction{}
	if err := session.Session.Query(`SELECT id, name, value, user_id FROM transaction WHERE id = ?`, id).Scan(&transaction.ID,
		&transaction.Name, &transaction.Value, &transaction.Value); err != nil {
		log.Println(err)
		return Transaction{}, err
	}
	return transaction, nil
}

func (session Session) GetAllTransactionsByUserId(userID int64) ([]Transaction, error) {
	var transactions []Transaction

	iter := session.Session.Query(`SELECT id, name, value, user_id FROM transaction WHERE user_id = ?`, userID).Iter()
	var transaction Transaction

	for iter.Scan(&transaction.ID, &transaction.Name, &transaction.Value, &transaction.UserID) {
		transactions = append(transactions, transaction)
	}

	if err := iter.Close(); err != nil {
		log.Println(err)
		return nil, err
	}

	return transactions, nil
}

func (session Session) SaveUser(id int64, firstName string, lastName string) {
	if err := session.Session.Query(`INSERT INTO user (id, first_name, last_name) VALUES (?, ?, ?)`,
		id, firstName, lastName).Exec(); err != nil {
		log.Println(err)
		return
	}
}

func (session Session) GetUser(id int64) (User, error) {
	user := User{}
	if err := session.Session.Query(`SELECT id, first_name, last_name FROM user WHERE id = ?`, id).Scan(&user.ID,
		&user.FirstName, &user.LastName); err != nil {
		log.Println(err)
		return User{}, err
	}
	return user, nil
}

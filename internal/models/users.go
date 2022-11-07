package models

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID       int
	Name     string
	Nickname string
	Balance  int
	Email    string
	Password string
	Created  time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name string, nickname string, balance int, email string, password string) (int, error) {
	stmt := `insert into users(name, nickname, balance, email, password, created)
	Values(?, ?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, name, nickname, balance, email, password, time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

func (m *UserModel) GetById(id int) (*User, error) {
	stmt := `SELECT * FROM users
	WHERE user_id = ?`

	row := m.DB.QueryRow(stmt, id)

	u := &User{}

	err := row.Scan(&u.ID, &u.Name, &u.Nickname, &u.Balance, &u.Email, &u.Password, &u.Created)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return u, nil
}

func (m *UserModel) CreateInitialUsers() error {
	name := "Zaur"
	nickname := "Lagmazavr"
	balance := 100
	email := "zaur@gmail.com"
	password := "password"
	_, err := m.Insert(name, nickname, balance, email, password)
	if err != nil {
		return err
	}

	return nil
}

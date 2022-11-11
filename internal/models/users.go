package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	ID       int
	Name     string
	Nickname string
	Balance  int
	Email    string
	Password []byte
	Created  time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name string, nickname string, balance int, email string, password string) error {
	stmt := `insert into users(name, nickname, balance, email, hashed_password, created)
	Values(?, ?, ?, ?, ?, UTC_TIMESTAMP)`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(stmt, name, nickname, balance, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	// Retrieve the id and hashed password associated with the given email. If
	// no matching email exists we return the ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte
	stmt := "SELECT user_id, hashed_password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
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

func (m *UserModel) UpdateUserInfo(id int, name, nickname string) error {
	stmt := `Update users
	set name=?, nickname=?
	where user_id=?`

	_, err := m.DB.Exec(stmt, name, nickname, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) UpdateBalance(id, balance int) error {
	stmt := `Update users
	set balance=?
	where user_id=?`

	_, err := m.DB.Exec(stmt, balance, id)
	if err != nil {
		return err
	}
	return nil
}
func (m *UserModel) CreateInitialUsers() error {
	name := "Zaur"
	nickname := "Lagmazavr"
	balance := 100
	email := "zaur@gmail.com"
	password := "password"

	err := m.Insert(name, nickname, balance, email, password)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) UpdatePassword(id int, password string) error {
	stmt := `Update users
	set hashed_password=?
	where user_id=?`

	_, err := m.DB.Exec(stmt, password, id)
	if err != nil {
		return err
	}
	return nil
}
func (m *UserModel) UserPurchasedGamesView(id int) ([]*Game, error) {
	stmt := `SELECT * FROM user_purchased_games 
	where user_purchased_games.game_id=?`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*Game

	for rows.Next() {
		game := &Game{}
		err := rows.Scan(&game.ID, &game.CreatedBy, &game.Name, &game.Description, &game.Cost, &game.ReleaseDate, &game.AuthorName)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return games, nil
}

func (m *UserModel) UserCreatedGamesView(id int) ([]*Game, error) {
	stmt := `SELECT g.game_id, g.created_by, g.name game_name, g.description, g.cost, g.release_year, u.name user_name
FROM games g JOIN users u ON g.created_by = u.user_id
where g.created_by=?;`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*Game

	for rows.Next() {
		game := &Game{}
		err := rows.Scan(&game.ID, &game.CreatedBy, &game.Name, &game.Description, &game.Cost, &game.ReleaseDate, &game.AuthorName)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return games, nil
}

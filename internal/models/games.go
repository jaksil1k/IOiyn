package models

import (
	"database/sql"
	"errors"
	"time"
)

type Game struct {
	ID          int
	CreatedBy   int
	AuthorName  string
	Name        string
	Description string
	Cost        int
	ReleaseDate time.Time
}

type GameModel struct {
	DB *sql.DB
}

func (m *GameModel) Insert(createdBy int, name string, description string, cost int, releaseYear time.Time) (int, error) {
	stmt := `insert into games(created_by, name, description, cost, release_year)
	Values(?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, createdBy, name, description, cost, releaseYear)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *GameModel) GetById(id int) (*Game, error) {
	stmt := `SELECT * FROM games
	WHERE game_id = ?`

	row := m.DB.QueryRow(stmt, id)

	g := &Game{}

	err := row.Scan(&g.ID, &g.CreatedBy, &g.Name, &g.Description, &g.Cost, &g.ReleaseDate)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return g, nil
}

func (m *GameModel) Latest() ([]*Game, error) {
	stmt := `SELECT * FROM games ORDER BY game_id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []*Game

	for rows.Next() {
		game := &Game{}
		err := rows.Scan(&game.ID, &game.CreatedBy, &game.Name, &game.Description, &game.Cost, &game.ReleaseDate)
		if err != nil {
			return nil, err
		}
		user, err := m.GetUserById(game.CreatedBy)
		if err != nil {
			return nil, err
		}
		game.AuthorName = user.Name
		games = append(games, game)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return games, nil
}

func (m *GameModel) GetUserById(id int) (*User, error) {
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

func (m *GameModel) CreateInitialGames() error {
	_, err := m.Insert(1, "dota", "kind of shit", 0, time.Date(2012, time.July, 9, 0, 0, 0, 0, time.UTC))
	if err != nil {
		return err
	}
	_, err = m.Insert(1, "cs go", "shit but not shit", 0, time.Date(2012, time.August, 21, 0, 0, 0, 0, time.UTC))
	if err != nil {
		return err
	}
	return nil
}

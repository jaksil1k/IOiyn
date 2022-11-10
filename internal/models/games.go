package models

import (
	"database/sql"
	"errors"
	"time"
)

type Game struct {
	ID          int
	CreatedBy   int
	Name        string
	Description string
	Cost        int
	ReleaseDate time.Time
	AuthorName  string
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
	stmt := `SELECT g.game_id, g.created_by, g.name , g.description, g.cost, g.release_year, u.name user_name
FROM games g JOIN users u ON g.created_by = u.user_id
WHERE g.game_id = ?;`

	row := m.DB.QueryRow(stmt, id)

	g := &Game{}

	err := row.Scan(&g.ID, &g.CreatedBy, &g.Name, &g.Description, &g.Cost, &g.ReleaseDate, &g.AuthorName)

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
	stmt := `SELECT * FROM games_user_name ORDER BY game_id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
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

func (m *GameModel) CreateInitialGames() error {
	_, err := m.Insert(1, "dota", "kind of shit", 0, time.Date(2012, time.July, 9, 0, 0, 0, 0, time.UTC))
	if err != nil {
		return err
	}
	_, err = m.Insert(1, "cs go", "simple гений", 0, time.Date(2012, time.August, 21, 0, 0, 0, 0, time.UTC))
	if err != nil {
		return err
	}
	_, err = m.Insert(1, "the last of us", "THE BEST GAME", 20, time.Date(2012, time.August, 21, 0, 0, 0, 0, time.UTC))
	if err != nil {
		return err
	}
	return nil
}

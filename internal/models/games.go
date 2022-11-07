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
	releaseYear time.Time
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

	err := row.Scan(&g.ID, &g.CreatedBy, &g.Name, &g.Description, &g.Cost, &g.releaseYear)

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
	return nil, nil
}

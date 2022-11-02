package models

import (
	"database/sql"
	"time"
)

type Game struct {
	ID          int
	GenreId     int
	Title       string
	Content     string
	Rating      float32
	releaseYear time.Time
}

type GameModel struct {
	DB *sql.DB
}

func (m *GameModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

func (m *GameModel) Get(id int) (*Game, error) {
	return nil, nil
}

func (m *GameModel) Latest() ([]*Game, error) {
	return nil, nil
}

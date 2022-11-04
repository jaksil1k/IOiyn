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

func (m *GameModel) CreateDB() error {
	createSchemaQuery := "create database if not exists games CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"
	useDBQuery := "USE game"
	createGamesQuery := `CREATE TABLE IF NOT EXISTS games(
	game_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    created_by INTEGER NOT NULL,
    genre_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
	description TEXT NOT NULL,
    rating NUMERIC(7, 2) NOT NULL,
    cost INTEGER NOT NULL,
    release_year int1 NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users(user_id),
    FOREIGN KEY (genre_id) REFERENCES genres(genre_id)
);`
	createUsersQuery := `CREATE TABLE IF NOT EXISTS users(
	user_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    nickname VARCHAR(30) NOT NULL,
	balance NUMERIC(7, 2) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password VARCHAR(100) NOT NULL
);`
	createGenreQuery := `CREATE TABLE IF NOT EXISTS genres(
	genre_id integer not null primary key auto_increment,
    genre_name varchar(30) NOT NULL
);`
	createPurchasedGamesQuery := `create table purchased_games(
	game_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    foreign key(game_id) references games(game_id),
    foreign key(user_id) references users(user_id)
);`
	//fmt.Println(createSchemaQuery, createUsersQuery, createGenreQuery, createGamesQuery, createPurchasedGamesQuery)
	_, err := m.DB.Exec(createSchemaQuery)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(useDBQuery)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(createUsersQuery)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(createGenreQuery)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(createGamesQuery)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(createPurchasedGamesQuery)
	if err != nil {
		return err
	}

	return nil
}

func (m *GameModel) InsertIntoGames(createdBy int, genreId int, name string, description string, rating float32, cost int, releaseYear int8) (int, error) {
	stmt := `insert into games(created_by, genre_id, name, description, rating, cost, release_year)
	Values(?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, createdBy, genreId, name, description, rating, cost, releaseYear)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *GameModel) InsertIntoUsers(name string, nickname string, balance float32, email string, password string) (int, error) {
	stmt := `insert into users(name, nickname, balance, email, password)
	Values(?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, name, nickname, balance, email, password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

	return 0, nil
}

func (m *GameModel) GetFromGames(id int) (*Game, error) {
	return nil, nil
}

func (m *GameModel) LatestFromGames() ([]*Game, error) {
	return nil, nil
}

func (m *GameModel) DropDB() error {
	dropPurchased := "DROP TABLE IF EXISTS purchased_games"
	dropGames := "DROP TABLE IF EXISTS games"
	dropUsers := "DROP TABLE IF EXISTS users"
	dropGenre := "DROP TABLE IF EXISTS genres"
	dropSchema := "DROP SCHEMA IF EXISTS game"

	_, err := m.DB.Exec(dropPurchased)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(dropGames)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(dropUsers)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(dropGenre)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(dropSchema)
	if err != nil {
		return err
	}

	return nil
}

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

type User struct {
	ID       int
	Name     string
	Nickname string
	Balance  int
	Email    int
	Password string
}

type GameModel struct {
	DB *sql.DB
}

func (m *GameModel) CreateTablesInDB() error {
	//createSchemaQuery := "create database if not exists game CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"
	//useDBQuery := "USE game"
	createGamesQuery := `CREATE TABLE IF NOT EXISTS games(
	game_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    created_by INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
	description TEXT NOT NULL,
    cost INTEGER NOT NULL,
    release_year DATE NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users(user_id)
);`
	createUsersQuery := `CREATE TABLE IF NOT EXISTS users(
	user_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    nickname VARCHAR(30) NOT NULL,
	balance INTEGER NOT NULL,
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
	createGamesGenresQuery := `create table games_genres(
	game_id INTEGER NOT NULL,
	genre_id INTEGER NOT NULL,
	FOREIGN KEY(game_id) REFERENCES games(game_id),
	FOREIGN KEY(genre_id) REFERENCES genres(genre_id)
);`
	//fmt.Println(createSchemaQuery, createUsersQuery, createGenreQuery, createGamesQuery, createPurchasedGamesQuery)
	//_, err := m.DB.Exec(createSchemaQuery)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = m.DB.Exec(useDBQuery)
	//if err != nil {
	//	return err
	//}

	_, err := m.DB.Exec(createUsersQuery)
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

	_, err = m.DB.Exec(createGamesGenresQuery)
	if err != nil {
		return err
	}

	return nil
}

func (m *GameModel) InsertIntoGames(createdBy int, name string, description string, cost int, releaseYear time.Time) (int, error) {
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

func (m *GameModel) InsertIntoUsers(name string, nickname string, balance int, email string, password string) (int, error) {
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
}

func (m *GameModel) GetByUserId(id int) (*User, error) {
	stmt := `SELECT * FROM users
	WHERE user_id = ?`

	row := m.DB.QueryRow(stmt, id)

	u := &User{}

	err := row.Scan(&u.ID, &u.Name, &u.Nickname, &u.Balance, &u.Email, &u.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return u, nil
}

func (m *GameModel) GetByGameId(id int) (*Game, error) {
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

func (m *GameModel) LatestFromGames() ([]*Game, error) {
	return nil, nil
}

func (m *GameModel) DropTablesInDB() error {
	dropPurchased := "DROP TABLE IF EXISTS purchased_games"
	dropGamesGenres := "DROP TABLE IF EXISTS games_genres"
	dropGames := "DROP TABLE IF EXISTS games"
	dropUsers := "DROP TABLE IF EXISTS users"
	dropGenre := "DROP TABLE IF EXISTS genres"
	//dropSchema := "DROP SCHEMA IF EXISTS game"

	_, err := m.DB.Exec(dropPurchased)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(dropGamesGenres)
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

	//_, err = m.DB.Exec(dropSchema)
	//if err != nil {
	//	return err
	//}

	return nil
}

package models

import "database/sql"

type DBModel struct {
	DB *sql.DB
}

func (m *DBModel) CreateTables() error {
	//createSchemaQuery := "create database if not exists game CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"
	//useDBQuery := "USE game"
	createGamesQuery := `CREATE TABLE IF NOT EXISTS games(
	game_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    created_by INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
	description TEXT NOT NULL,
    cost INTEGER NOT NULL,
    release_year DATE NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users(user_id)
);`
	createUsersQuery := `CREATE TABLE IF NOT EXISTS users(
	user_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    nickname VARCHAR(255) NOT NULL,
	balance INTEGER NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
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
	createSession := `CREATE TABLE sessions (
		token CHAR(43) PRIMARY KEY,
		data BLOB NOT NULL,
		expiry TIMESTAMP(6) NOT NULL
	);`

	userEmailConstraint := `ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);`

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

	_, err = m.DB.Exec(createSession)
	if err != nil {
		return err
	}

	createSession = `CREATE INDEX sessions_expiry_idx ON sessions (expiry);`
	_, err = m.DB.Exec(createSession)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(userEmailConstraint)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) DropTables() error {
	dropPurchased := "DROP TABLE IF EXISTS purchased_games"
	dropGamesGenres := "DROP TABLE IF EXISTS games_genres"
	dropGames := "DROP TABLE IF EXISTS games"
	dropUsers := "DROP TABLE IF EXISTS users"
	dropGenre := "DROP TABLE IF EXISTS genres"
	dropSession := "DROP TABLE IF EXISTS sessions"
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

	_, err = m.DB.Exec(dropSession)
	if err != nil {
		return err
	}
	//_, err = m.DB.Exec(dropSchema)
	//if err != nil {
	//	return err
	//}

	return nil
}

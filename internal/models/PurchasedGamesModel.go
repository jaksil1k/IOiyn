package models

import (
	"database/sql"
)

type PurchasedGamesStruct struct {
}

type PurchasedGamesModel struct {
	DB *sql.DB
}

func (m *PurchasedGamesModel) Insert(gameId, userId int) error {
	stmt := `insert into purchased_games(game_id, user_id)
	Values(?, ?)`

	_, err := m.DB.Exec(stmt, gameId, userId)

	if err != nil {
		return err
	}

	return nil
}

func (m *PurchasedGamesModel) IsExists(gameId, userId int) (bool, error) {
	stmt := `SELECT EXISTS(SELECT * FROM purchased_games WHERE game_id=? AND user_id =?)`

	row := m.DB.QueryRow(stmt, gameId, userId)

	var isOk bool

	err := row.Scan(&isOk)

	if err != nil {
		return false, err
	}

	return isOk, nil
}

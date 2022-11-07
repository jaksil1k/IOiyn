package main

import "IOiyn.kz/internal/models"

type templateData struct {
	Game  *models.Game
	Games []*models.Game
	User  *models.User
}

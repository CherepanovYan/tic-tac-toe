package domain

import (
	"tic-tac-toe/internal/domain/model"

	"github.com/google/uuid"
)

type Repository interface {
	SaveGame(game *model.Game) error
	GetGameByUUID(id uuid.UUID) (*model.Game, error)
	GetCompletedGamesByUserUUID(userUUID uuid.UUID) ([]*model.Game, error)
	GetLeaderboard(limit int) ([]*model.LeaderboardEntry, error)
	CreateUser(user *model.User) error
	GetUserByLogin(login string) (*model.User, error)
	GetAvailableGames() ([]*model.Game, error)
	GetUserByUUID(id uuid.UUID) (*model.User, error)
}

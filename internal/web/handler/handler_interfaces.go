package handler

import (
	"tic-tac-toe/internal/domain/model"
	web "tic-tac-toe/internal/web/model"

	"github.com/google/uuid"
)

type GameService interface {
	CreatePlayerGame(userX uuid.UUID) *model.Game
	CreateComputerGame(userX uuid.UUID) *model.Game

	JoinGame(game *model.Game, userID uuid.UUID) error

	GetAvailableGames() ([]*model.Game, error)
	GetCompletedGamesByUserUUID(userUUID uuid.UUID) ([]*model.Game, error)
	GetLeaderboard(limit int) ([]*model.LeaderboardEntry, error)

	MakeMove(
		game *model.Game,
		row int,
		column int,
		player uuid.UUID,
	) error

	GetGameInfo(id uuid.UUID) (*model.Game, error)

	MinimaxMove(game *model.Game) error

	ValidateBoardUnchanged(
		old,
		new model.PlayingField,
	) error
}

type WebMappers interface {
	GameToResponse(game *model.Game) *web.GameResponse
}

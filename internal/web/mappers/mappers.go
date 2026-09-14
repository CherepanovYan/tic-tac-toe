package mappers

import (
	"time"

	"tic-tac-toe/internal/domain/model"
	web "tic-tac-toe/internal/web/model"
)

type WebMappers struct{}

func NewWebMappers() *WebMappers {
	return &WebMappers{}
}

func (m *WebMappers) GameToResponse(game *model.Game) *web.GameResponse {
	if game == nil {
		return nil
	}

	return &web.GameResponse{
		UUID:          game.ID.String(),
		PlayingField:  game.PlayingField,
		GameType:      string(game.GameType),
		Status:        string(game.Status),
		UserX:         game.UserX.String(),
		UserO:         game.UserO.String(),
		CurrentPlayer: game.CurrentPlayer.String(),
		CreatedAt:     game.CreatedAt.Format(time.RFC3339),
	}
}

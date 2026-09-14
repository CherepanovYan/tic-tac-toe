package model

import (
	"tic-tac-toe/internal/domain/model"
	"time"

	"github.com/google/uuid"
)

type BoardDTO [model.SizeField][model.SizeField]int

type GameDTO struct {
	ID            uuid.UUID `db:"id"`
	BoardDTO      BoardDTO  `db:"board"`
	StatusDTO     string    `db:"status"`
	GameTypeDTO   string    `db:"game_type"`
	UserX         uuid.UUID `db:"userX"`
	UserO         uuid.UUID `db:"userO"`
	CurrentPlayer uuid.UUID `db:"currentPlayer"`
	CreatedAt     time.Time `db:"created_at"`
}

type UserDTO struct {
	ID       uuid.UUID `db:"id"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
}

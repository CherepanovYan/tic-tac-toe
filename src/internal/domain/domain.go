package domain

import (
	"tic-tac-toe/internal/domain/model"

	"github.com/google/uuid"
)

type Repository interface {
}

type domain struct {
	repo Repository
}

func NewDomain(repository Repository) *domain {
	return &domain{repo: repository}
}

func (d *domain) GetNextMove() (row, col int, err error)       { return 1, 1, nil } // minimax
func (d *domain) ValidateBoardUnchanged() error                { return nil }
func (d *domain) IsGameFinished() (finished bool, winner int8) { return false, 1 }
func (d *domain) NewGame() model.Game {
	return model.Game{
		UUID:         uuid.New(),
		PlayingField: model.NewPlayingField(),
	}
}

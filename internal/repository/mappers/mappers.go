package mappers

import (
	"tic-tac-toe/internal/domain/model"
	data "tic-tac-toe/internal/repository/model"
)

func ToGameDTO(game *model.Game) *data.GameDTO {
	if game == nil {
		return nil
	}

	var board data.BoardDTO

	for i := 0; i < model.SizeField; i++ {
		for j := 0; j < model.SizeField; j++ {
			board[i][j] = game.PlayingField[i][j]
		}
	}

	return &data.GameDTO{
		ID:            game.ID,
		BoardDTO:      board,
		StatusDTO:     string(game.Status),
		GameTypeDTO:   string(game.GameType),
		UserX:         game.UserX,
		UserO:         game.UserO,
		CurrentPlayer: game.CurrentPlayer,
		CreatedAt:     game.CreatedAt,
	}
}

func ToGame(gameDTO *data.GameDTO) *model.Game {
	if gameDTO == nil {
		return nil
	}

	var board model.PlayingField

	for i := 0; i < model.SizeField; i++ {
		for j := 0; j < model.SizeField; j++ {
			board[i][j] = gameDTO.BoardDTO[i][j]
		}
	}

	return &model.Game{
		ID:            gameDTO.ID,
		PlayingField:  board,
		Status:        model.GameStatus(gameDTO.StatusDTO),
		GameType:      model.GameType(gameDTO.GameTypeDTO),
		UserX:         gameDTO.UserX,
		UserO:         gameDTO.UserO,
		CurrentPlayer: gameDTO.CurrentPlayer,
		CreatedAt:     gameDTO.CreatedAt,
	}
}

func ToUser(userDTO *data.UserDTO) *model.User {
	if userDTO == nil {
		return nil
	}
	return &model.User{
		UUID:     userDTO.ID,
		Login:    userDTO.Login,
		Password: userDTO.Password,
	}
}

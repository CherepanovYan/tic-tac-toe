package domain

import (
	"tic-tac-toe/internal/domain/model"

	"github.com/google/uuid"
)

type gameService struct {
	repo Repository
}

func NewGameService(repository Repository) *gameService {
	return &gameService{
		repo: repository,
	}
}

func (d *gameService) CreateComputerGame(userX uuid.UUID) *model.Game {
	game := model.NewComputerGame(userX)

	if err := d.repo.SaveGame(game); err != nil {
		return nil
	}

	return game
}

func (d *gameService) CreatePlayerGame(userX uuid.UUID) *model.Game {
	game := model.NewPlayerGame(userX)

	if err := d.repo.SaveGame(game); err != nil {
		return nil
	}

	return game
}

func (d *gameService) JoinGame(game *model.Game, userID uuid.UUID) error {
	if game == nil {
		return model.ErrorNotFound
	}

	if userID == uuid.Nil {
		return model.ErrorUnauthorized
	}

	if game.GameType != model.GameTypePlayer {
		return model.ErrorGameNotInProgress
	}

	if game.Status != model.StatusWaiting {
		return model.ErrorGameNotInProgress
	}

	if game.UserX == userID {
		return model.ErrorInvalidPlayer
	}

	game.UserO = userID
	game.Status = model.StatusXTurn
	game.CurrentPlayer = game.UserX

	return d.repo.SaveGame(game)
}

func (d *gameService) MakeMove(
	game *model.Game,
	row int,
	column int,
	userID uuid.UUID,
) error {
	if game == nil {
		return model.ErrorNotFound
	}

	if game.Status == model.StatusWaiting {
		return model.ErrorGameNotInProgress
	}

	if game.Status != model.StatusXTurn &&
		game.Status != model.StatusOTurn {
		return model.ErrorGameNotInProgress
	}

	if userID == uuid.Nil {
		return model.ErrorUnauthorized
	}

	if userID != game.CurrentPlayer {
		return model.ErrorNotYourTurn
	}

	var player int

	switch userID {
	case game.UserX:
		player = model.PlayerX

	case game.UserO:
		player = model.PlayerO

	default:
		return model.ErrorInvalidPlayer
	}

	if !game.PlayingField.Set(row, column, player) {
		return model.ErrorInvalidMove
	}

	d.UpdateGameStatus(game)

	if game.Status == model.StatusXWon ||
		game.Status == model.StatusOWon ||
		game.Status == model.StatusDraw {
		return nil
	}

	if game.GameType == model.GameTypeComputer {
		game.CurrentPlayer = uuid.Nil
		game.Status = model.StatusOTurn

		return d.repo.SaveGame(game)
	}

	switch game.CurrentPlayer {
	case game.UserX:
		game.CurrentPlayer = game.UserO
		game.Status = model.StatusOTurn

	case game.UserO:
		game.CurrentPlayer = game.UserX
		game.Status = model.StatusXTurn

	default:
		return model.ErrorInvalidPlayer
	}

	return d.repo.SaveGame(game)
}

func (d *gameService) GetGameInfo(id uuid.UUID) (*model.Game, error) {
	if id == uuid.Nil {
		return nil, model.ErrorInvalidUUID
	}

	game, err := d.repo.GetGameByUUID(id)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (d *gameService) MinimaxMove(game *model.Game) error {
	if game == nil {
		return model.ErrorNotFound
	}

	if game.GameType != model.GameTypeComputer {
		return model.ErrorInvalidPlayer
	}

	if game.Status != model.StatusOTurn {
		return model.ErrorNotYourTurn
	}

	bestScore := -1000
	bestRow := -1
	bestCol := -1

	for i := 0; i < model.SizeField; i++ {
		for j := 0; j < model.SizeField; j++ {
			if game.PlayingField[i][j] != model.Empty {
				continue
			}

			game.PlayingField[i][j] = model.PlayerO

			score := d.minimax(game, 0, false)

			game.PlayingField[i][j] = model.Empty

			if score > bestScore {
				bestScore = score
				bestRow = i
				bestCol = j
			}
		}
	}

	if bestRow == -1 || bestCol == -1 {
		return model.ErrorInvalidMove
	}

	game.PlayingField[bestRow][bestCol] = model.PlayerO

	d.UpdateGameStatus(game)

	if game.Status == model.StatusOWon ||
		game.Status == model.StatusDraw {
		return nil
	}

	game.CurrentPlayer = game.UserX
	game.Status = model.StatusXTurn

	return d.repo.SaveGame(game)
}

func (d *gameService) minimax(
	game *model.Game,
	depth int,
	isMaximizing bool,
) int {
	if game == nil {
		return 0
	}

	status := d.getGameStatus(game)

	switch status {
	case model.StatusOWon:
		return 10 - depth

	case model.StatusXWon:
		return depth - 10

	case model.StatusDraw:
		return 0
	}

	eval := 100

	factor := -1
	player := model.PlayerO

	if !isMaximizing {
		factor = 1
		player = model.PlayerX
	}

	eval *= factor

	for i := 0; i < model.SizeField; i++ {
		for j := 0; j < model.SizeField; j++ {

			if game.PlayingField[i][j] == model.Empty {
				game.PlayingField[i][j] = player

				tmpEval := d.minimax(
					game,
					depth+1,
					!isMaximizing,
				)

				game.PlayingField[i][j] = model.Empty

				if isMaximizing {
					if tmpEval > eval {
						eval = tmpEval
					}
				} else {
					if tmpEval < eval {
						eval = tmpEval
					}
				}
			}
		}
	}

	return eval
}

func (d *gameService) ValidateBoardUnchanged(
	old,
	new model.PlayingField,
) error {
	if len(old) != model.SizeField ||
		len(new) != model.SizeField {
		return model.ErrorInvalidBoard
	}

	for i := 0; i < model.SizeField; i++ {
		for j := 0; j < model.SizeField; j++ {

			if old[i][j] != model.Empty &&
				old[i][j] != new[i][j] {
				return model.ErrorInvalidMove
			}
		}
	}

	return nil
}

func (d *gameService) UpdateGameStatus(game *model.Game) {
	if game == nil {
		return
	}

	if game.GameType == model.GameTypePlayer &&
		game.UserO == uuid.Nil {
		game.Status = model.StatusWaiting

		_ = d.repo.SaveGame(game)
		return
	}

	status := d.getGameStatus(game)

	switch status {
	case model.StatusXWon:
		game.Status = model.StatusXWon
		_ = d.repo.SaveGame(game)
		return

	case model.StatusOWon:
		game.Status = model.StatusOWon
		_ = d.repo.SaveGame(game)
		return

	case model.StatusDraw:
		game.Status = model.StatusDraw
		_ = d.repo.SaveGame(game)
		return
	}

	switch game.CurrentPlayer {
	case game.UserX:
		game.Status = model.StatusXTurn

	case game.UserO:
		game.Status = model.StatusOTurn

	default:
		if game.GameType == model.GameTypeComputer {
			game.Status = model.StatusOTurn
		}
	}

	_ = d.repo.SaveGame(game)
}

func (d *gameService) getGameStatus(
	game *model.Game,
) model.GameStatus {
	if game == nil {
		return model.StatusUnknown
	}

	n := model.SizeField
	board := game.PlayingField

	dirs := [][2]int{
		{0, 1},
		{1, 0},
		{1, 1},
		{1, -1},
	}

	// Проверяем победу.
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			player := board[i][j]

			if player == model.Empty {
				continue
			}

			for _, dir := range dirs {
				count := 1

				x := i
				y := j

				for count < n {
					x += dir[0]
					y += dir[1]

					if x < 0 ||
						x >= n ||
						y < 0 ||
						y >= n {
						break
					}

					if board[x][y] != player {
						break
					}

					count++
				}

				if count == n {
					switch player {
					case model.PlayerX:
						return model.StatusXWon

					case model.PlayerO:
						return model.StatusOWon
					}
				}
			}
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if board[i][j] == model.Empty {
				return model.StatusUnknown
			}
		}
	}

	return model.StatusDraw
}

func (d *gameService) GetAvailableGames() ([]*model.Game, error) {
	return d.repo.GetAvailableGames()
}

func (d *gameService) GetCompletedGamesByUserUUID(
	userUUID uuid.UUID,
) ([]*model.Game, error) {
	if userUUID == uuid.Nil {
		return nil, model.ErrorInvalidUUID
	}

	return d.repo.GetCompletedGamesByUserUUID(userUUID)
}

func (d *gameService) GetLeaderboard(
	limit int,
) ([]*model.LeaderboardEntry, error) {
	return d.repo.GetLeaderboard(limit)
}

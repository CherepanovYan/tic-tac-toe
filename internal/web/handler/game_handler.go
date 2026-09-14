package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"tic-tac-toe/internal/domain"
	"tic-tac-toe/internal/domain/model"
	web "tic-tac-toe/internal/web"
	webModel "tic-tac-toe/internal/web/model"
)

type GameHandler struct {
	service    GameService
	repo       domain.Repository
	webMappers WebMappers
}

func NewGameHandler(
	service GameService,
	repo domain.Repository,
	webMappers WebMappers,
) *GameHandler {
	return &GameHandler{
		service:    service,
		repo:       repo,
		webMappers: webMappers,
	}
}

func (g *GameHandler) CreateNewGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		webModel.ResponseError(
			w,
			http.StatusMethodNotAllowed,
			model.ErrorInvalidMethod,
			"Only GET method is allowed",
		)
		return
	}

	user, ok := web.UserFromContext(r.Context())
	if !ok {
		webModel.ResponseError(
			w,
			http.StatusUnauthorized,
			model.ErrorUnauthorized,
			"user not found",
		)
		return
	}

	gameType := r.URL.Query().Get("type")

	var game *model.Game

	switch gameType {
	case "", "computer":
		game = g.service.CreateComputerGame(user.UUID)

	case "player":
		game = g.service.CreatePlayerGame(user.UUID)

	default:
		webModel.ResponseError(
			w,
			http.StatusBadRequest,
			model.ErrorInvalidPlayer,
			"Invalid game type. Use 'player' or 'computer'",
		)
		return
	}

	if game == nil {
		webModel.ResponseError(
			w,
			http.StatusInternalServerError,
			model.ErrorNotFound,
			"Failed to create a new game",
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(map[string]string{
		"id": game.ID.String(),
	})

	if err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}
}

func (g *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		webModel.ResponseError(
			w,
			http.StatusMethodNotAllowed,
			model.ErrorInvalidMethod,
			"Only POST method is allowed",
		)
		return
	}

	user, ok := web.UserFromContext(r.Context())
	if !ok {
		webModel.ResponseError(
			w,
			http.StatusUnauthorized,
			model.ErrorUnauthorized,
			"user not found",
		)
		return
	}

	idValue := r.PathValue("UUID")
	id, err := uuid.Parse(idValue)
	if err != nil {
		webModel.ResponseError(
			w,
			http.StatusBadRequest,
			model.ErrorInvalidUUID,
			"Invalid user UUID",
		)
		return
	}

	game, err := g.service.GetGameInfo(id)
	if err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}

	err = g.service.JoinGame(game, user.UUID)
	if err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}

	response := g.webMappers.GameToResponse(game)

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)

	err = encoder.Encode(response)
	if err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}
}

func (g *GameHandler) GameInfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodGet {
		webModel.ResponseError(
			w,
			http.StatusMethodNotAllowed,
			model.ErrorInvalidMethod,
			"Only GET method is allowed",
		)
		return
	}

	idValue := r.PathValue("UUID")
	id, err := uuid.Parse(idValue)
	if err != nil {
		webModel.ResponseError(
			w,
			http.StatusBadRequest,
			model.ErrorInvalidUUID,
			"Invalid user UUID",
		)
		return
	}

	game, err := g.service.GetGameInfo(id)
	if err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}

	response := g.webMappers.GameToResponse(game)

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)

	err = encoder.Encode(response)
	if err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}
}

func (g *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		webModel.ResponseError(
			w,
			http.StatusMethodNotAllowed,
			model.ErrorInvalidMethod,
			"Only POST method is allowed",
		)
		return
	}

	currentPlayer, ok := web.UserFromContext(r.Context())
	if !ok {
		webModel.ResponseError(
			w,
			http.StatusUnauthorized,
			model.ErrorUnauthorized,
			"user not found",
		)
		return
	}

	var move webModel.MoveRequest

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&move)
	if err != nil {
		webModel.ResponseError(
			w,
			http.StatusBadRequest,
			err,
			"Invalid request",
		)
		return
	}

	idValue := r.PathValue("UUID")
	id, err := uuid.Parse(idValue)
	if err != nil {
		webModel.ResponseError(
			w,
			http.StatusBadRequest,
			model.ErrorInvalidUUID,
			"Invalid user UUID",
		)
		return
	}

	game, err := g.service.GetGameInfo(id)
	if err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}

	oldPlayingField := game.PlayingField.Copy()

	err = g.service.MakeMove(
		game,
		move.Row,
		move.Column,
		currentPlayer.UUID,
	)

	if err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}

	if game.GameType == model.GameTypeComputer &&
		game.Status == model.StatusOTurn {

		err = g.service.MinimaxMove(game)
		if err != nil {
			code, message := webModel.ErrorConvertor(err)

			webModel.ResponseError(
				w,
				code,
				err,
				message,
			)
			return
		}
	}

	err = g.service.ValidateBoardUnchanged(
		oldPlayingField,
		game.PlayingField,
	)

	if err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}

	response := g.webMappers.GameToResponse(game)

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)

	err = encoder.Encode(response)
	if err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}
}

func (g *GameHandler) GetAvailableGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		webModel.ResponseError(
			w,
			http.StatusMethodNotAllowed,
			model.ErrorInvalidMethod,
			"Only GET method is allowed",
		)
		return
	}

	games, err := g.service.GetAvailableGames()
	if err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}

	responses := make([]*webModel.GameResponse, 0, len(games))

	for _, game := range games {
		responses = append(
			responses,
			g.webMappers.GameToResponse(game),
		)
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(responses); err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}
}

func (g *GameHandler) GetCompletedGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		webModel.ResponseError(
			w,
			http.StatusMethodNotAllowed,
			model.ErrorInvalidMethod,
			"Only GET method is allowed",
		)
		return
	}

	user, ok := web.UserFromContext(r.Context())
	if !ok {
		webModel.ResponseError(
			w,
			http.StatusUnauthorized,
			model.ErrorUnauthorized,
			"user not found",
		)
		return
	}

	games, err := g.service.GetCompletedGamesByUserUUID(user.UUID)
	if err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}

	responses := make([]*webModel.GameResponse, 0, len(games))
	for _, game := range games {
		responses = append(responses, g.webMappers.GameToResponse(game))
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(responses); err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}
}

func (g *GameHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		webModel.ResponseError(
			w,
			http.StatusMethodNotAllowed,
			model.ErrorInvalidMethod,
			"Only GET method is allowed",
		)
		return
	}

	limit := 10
	queryLimit := r.URL.Query().Get("limit")
	if queryLimit != "" {
		parsedLimit, err := strconv.Atoi(queryLimit)
		if err != nil || parsedLimit <= 0 {
			webModel.ResponseError(
				w,
				http.StatusBadRequest,
				model.ErrorInvalidMove,
				"limit must be a positive integer",
			)
			return
		}
		limit = parsedLimit
	}

	entries, err := g.service.GetLeaderboard(limit)
	if err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}

	responses := make([]*webModel.LeaderboardResponse, 0, len(entries))
	for _, entry := range entries {
		responses = append(responses, &webModel.LeaderboardResponse{
			UUID:     entry.UUID.String(),
			Login:    entry.Login,
			WinRatio: entry.WinRatio,
		})
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(responses); err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
		return
	}
}

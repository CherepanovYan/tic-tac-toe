package game

import (
	"encoding/json"

	"net/http"
	"tic-tac-toe/internal/domain/model"
)

type GameService interface {
	GetNextMove() (row, col int, err error) // minimax
	ValidateBoardUnchanged() error
	IsGameFinished() (finished bool, winner int8)
}

type NewGame interface {
	NewGame() model.Game
}

type gameHandler struct {
	service GameService
	newGame NewGame
}

func NewGameHendler(service GameService, newGame NewGame) *gameHandler {
	return &gameHandler{service: service,
		newGame: newGame}
}

func (g *gameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//uuid := r.URL.Query().Get("current_game_UUID")

	decoder := json.NewDecoder(r.Body)
	
	var game model.Game
	err := decoder.Decode(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game.PlayingField.Field[1][1] = 1 // player move
	encoder := json.NewEncoder(w)
	err = encoder.Encode(game.PlayingField)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (g *gameHandler) CreateNewGame(w http.ResponseWriter, r *http.Request) {
	game := g.newGame.NewGame()

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

package handler

import (
	"net/http"
)

type gameHendler interface {
	CreateNewGame(w http.ResponseWriter, r *http.Request)
	MakeMove(w http.ResponseWriter, r *http.Request)
}

func RegisterRouter(mux *http.ServeMux, gameHendler gameHendler) {
	mux.HandleFunc("GET /newGame", gameHendler.CreateNewGame)
	mux.HandleFunc("POST /game/{current_game_UUID}", gameHendler.MakeMove)
}

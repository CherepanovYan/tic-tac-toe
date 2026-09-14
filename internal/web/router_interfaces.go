package handler

import (
	"net/http"
)

type gameHandler interface {
	CreateNewGame(w http.ResponseWriter, r *http.Request)
	MakeMove(w http.ResponseWriter, r *http.Request)
	GameInfo(w http.ResponseWriter, r *http.Request)
	JoinGame(w http.ResponseWriter, r *http.Request)
	GetAvailableGames(w http.ResponseWriter, r *http.Request)
	GetCompletedGames(w http.ResponseWriter, r *http.Request)
	GetLeaderboard(w http.ResponseWriter, r *http.Request)
}

type authHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	RefreshAccessToken(w http.ResponseWriter, r *http.Request)
	RefreshRefreshToken(w http.ResponseWriter, r *http.Request)
}

type userHandler interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	GetCurrentUser(w http.ResponseWriter, r *http.Request)
}

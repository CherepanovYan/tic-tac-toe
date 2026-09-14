package model

import "tic-tac-toe/internal/domain/model"

type GameResponse struct {
	UUID          string             `json:"uuid"`
	PlayingField  model.PlayingField `json:"playing_field"`
	GameType      string             `json:"game_type"`
	Status        string             `json:"status"`
	UserX         string             `json:"user_x"`
	UserO         string             `json:"user_o"`
	CurrentPlayer string             `json:"current_player"`
	CreatedAt     string             `json:"created_at"`
}

type MoveRequest struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

type LeaderboardResponse struct {
	UUID     string  `json:"uuid"`
	Login    string  `json:"login"`
	WinRatio float64 `json:"win_ratio"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

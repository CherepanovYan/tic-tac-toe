package model

import (
	"time"

	"github.com/google/uuid"
)

type GameType string

const (
	GameTypeComputer GameType = "COMPUTER"
	GameTypePlayer   GameType = "PLAYER"
)

type GameStatus string

const (
	StatusWaiting GameStatus = "WAITING"
	StatusXTurn   GameStatus = "X_TURN"
	StatusOTurn   GameStatus = "O_TURN"
	StatusXWon    GameStatus = "X_WON"
	StatusOWon    GameStatus = "O_WON"
	StatusDraw    GameStatus = "DRAW"
	StatusUnknown GameStatus = "UNKNOWN"
)

type Game struct {
	ID            uuid.UUID    `json:"id"`
	PlayingField  PlayingField `json:"playing_field"`
	GameType      GameType     `json:"game_type"`
	Status        GameStatus   `json:"status"`
	UserX         uuid.UUID    `json:"user_x"`
	UserO         uuid.UUID    `json:"user_o"`
	CurrentPlayer uuid.UUID    `json:"current_player"`
	CreatedAt     time.Time    `json:"created_at"`
}

func NewComputerGame(userX uuid.UUID) *Game {
	return &Game{
		ID:            uuid.New(),
		PlayingField:  NewPlayingField(),
		GameType:      GameTypeComputer,
		Status:        StatusXTurn,
		UserX:         userX,
		UserO:         uuid.Nil,
		CurrentPlayer: userX,
		CreatedAt:     time.Now().UTC(),
	}
}

func NewPlayerGame(userX uuid.UUID) *Game {
	return &Game{
		ID:            uuid.New(),
		PlayingField:  NewPlayingField(),
		GameType:      GameTypePlayer,
		Status:        StatusWaiting,
		UserX:         userX,
		UserO:         uuid.Nil,
		CurrentPlayer: uuid.Nil,
		CreatedAt:     time.Now().UTC(),
	}
}

type LeaderboardEntry struct {
	UUID     uuid.UUID `json:"uuid"`
	Login    string    `json:"login"`
	WinRatio float64   `json:"win_ratio"`
}

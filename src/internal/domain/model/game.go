package model

import ("github.com/google/uuid")

type Game struct {
	UUID         uuid.UUID `json:"uuid"`
	PlayingField *playingField `json:"playing_field"`
}
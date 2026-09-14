package model

import "github.com/google/uuid"

type User struct {
	UUID     uuid.UUID `json:"id"`
	Login    string    `json:"login"`
	Password string    `json:"-"`
}

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

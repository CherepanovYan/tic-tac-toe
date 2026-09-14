package model

import "errors"

var (
	ErrorInvalidMove          = errors.New("Invalid move")
	ErrorNotFound             = errors.New("Game not found")
	ErrorInvalidBoard         = errors.New("Invalid board")
	ErrorInvalidUUID          = errors.New("Invalid UUID format")
	ErrorInvalidMethod        = errors.New("Invalid method")
	ErrorGameNotInProgress    = errors.New("Game is not in progress")
	ErrorNotYourTurn          = errors.New("Not your turn")
	ErrorInvalidPlayer        = errors.New("Invalid player")
	ErrorUnauthorized         = errors.New("Unauthorized")
	ErrorInvalidCredentials   = errors.New("invalid credentials")
	ErrorLoginRequired        = errors.New("login is required")
	ErrorPasswordRequired     = errors.New("password is required")
	ErrorUserAlreadyExists    = errors.New("user already exists")
	ErrorRefreshTokenRequired = errors.New("refresh token required")
	ErrorGenerateAccessToken  = errors.New("error generate access token")
	ErrorGenerateRefreshToken = errors.New("error generate refresh token")
)

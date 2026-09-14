package model

import (
	"encoding/json"
	"net/http"
	"tic-tac-toe/internal/domain/model"
)

func ResponseError(w http.ResponseWriter, code int, error error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := ErrorResponse{
		Error:   error.Error(),
		Message: message,
		Code:    code,
	}

	json.NewEncoder(w).Encode(err)
}

func ErrorConvertor(err error) (int, string) {
	switch err {
	case model.ErrorNotFound:
		return http.StatusNotFound, "Game not found"
	case model.ErrorGameNotInProgress:
		return http.StatusBadRequest, "Game not in progress"
	case model.ErrorInvalidMove:
		return http.StatusBadRequest, "Invalid move"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}

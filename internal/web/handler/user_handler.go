package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"tic-tac-toe/internal/domain"
	"tic-tac-toe/internal/domain/model"
	web "tic-tac-toe/internal/web"
	webModel "tic-tac-toe/internal/web/model"
)

type UserHandler struct {
	service *domain.UserService
}

func NewUserHandler(service *domain.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.service.GetUserByUUID(id)
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

	response := map[string]interface{}{
		"id":    user.UUID,
		"login": user.Login,
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(response); err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
	}
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
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

	response := map[string]interface{}{
		"id":    user.UUID,
		"login": user.Login,
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(response); err != nil {
		code, message := webModel.ErrorConvertor(err)

		webModel.ResponseError(
			w,
			code,
			err,
			message,
		)
	}
}

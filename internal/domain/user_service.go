package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"tic-tac-toe/internal/domain/model"
)

type UserService struct {
	repository Repository
}

func NewUserService(repository Repository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) Register(request model.SignUpRequest) error {
	if request.Login == "" {
		return model.ErrorLoginRequired
	}

	if request.Password == "" {
		return model.ErrorPasswordRequired
	}

	_, err := s.repository.GetUserByLogin(request.Login)
	if err == nil {
		return model.ErrorUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		UUID:     uuid.New(),
		Login:    request.Login,
		Password: string(passwordHash),
	}

	if err := s.repository.CreateUser(user); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *UserService) Authenticate(login, password string) (*model.User, error) {
	user, err := s.repository.GetUserByLogin(login)
	if err != nil {
		return nil, model.ErrorInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return nil, model.ErrorInvalidCredentials
	}

	return user, nil
}

func (s *UserService) GetUserByUUID(id uuid.UUID) (*model.User, error) {
	if id == uuid.Nil {
		return nil, model.ErrorInvalidUUID
	}

	return s.repository.GetUserByUUID(id)
}

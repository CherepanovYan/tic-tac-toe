package domain

import (
	"tic-tac-toe/internal/domain/model"
)

type AuthService struct {
	userService *UserService
	jwtService  *JwtService
}

func NewAuthService(userService *UserService, jwtService *JwtService) *AuthService {
	return &AuthService{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (s *AuthService) Register(request model.SignUpRequest) error {
	return s.userService.Register(request)
}

func (s *AuthService) Authenticate(
	request model.JwtRequest,
) (*model.JwtResponse, error) {
	if request.Login == "" || request.Password == "" {
		return nil, model.ErrorInvalidCredentials
	}

	user, err := s.userService.Authenticate(
		request.Login,
		request.Password,
	)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.UUID)
	if err != nil {
		return nil, model.ErrorGenerateAccessToken
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.UUID)
	if err != nil {
		return nil, model.ErrorGenerateRefreshToken
	}

	return &model.JwtResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshAccessToken(
	request model.RefreshJwtRequest,
) (*model.JwtResponse, error) {
	if request.RefreshToken == "" {
		return nil, model.ErrorRefreshTokenRequired
	}

	if err := s.jwtService.ValidateRefreshToken(request.RefreshToken); err != nil {
		return nil, err
	}

	userUUID, err := s.jwtService.GetUUID(request.RefreshToken)
	if err != nil {
		return nil, err
	}

	_, err = s.userService.GetUserByUUID(userUUID)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtService.GenerateAccessToken(userUUID)
	if err != nil {
		return nil, model.ErrorGenerateAccessToken
	}

	return &model.JwtResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: request.RefreshToken,
	}, nil
}

func (s *AuthService) RefreshRefreshToken(
	request model.RefreshJwtRequest,
) (*model.JwtResponse, error) {
	if request.RefreshToken == "" {
		return nil, model.ErrorRefreshTokenRequired
	}

	if err := s.jwtService.ValidateRefreshToken(request.RefreshToken); err != nil {
		return nil, err
	}

	userUUID, err := s.jwtService.GetUUID(request.RefreshToken)
	if err != nil {
		return nil, err
	}

	_, err = s.userService.GetUserByUUID(userUUID)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtService.GenerateAccessToken(userUUID)
	if err != nil {
		return nil, model.ErrorGenerateAccessToken
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(userUUID)
	if err != nil {
		return nil, model.ErrorGenerateRefreshToken
	}

	return &model.JwtResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

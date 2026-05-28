package services

import (
	"errors"
	"time"

	"user-management-api/auth"
	"user-management-api/models"
	"user-management-api/repositories"
	"user-management-api/utils"

	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already registered")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService interface {
	Login(req models.LoginRequest) (models.AuthTokenResponse, error)
	Register(req models.RegisterRequest) (models.AuthTokenResponse, error)
}

type authService struct {
	userRepo   repositories.UserRepository
	logService LogService
	jwtSecret  string
	jwtExpiry  int
}

func NewAuthService(userRepo repositories.UserRepository, logService LogService, jwtSecret string, jwtExpiry int) AuthService {
	return &authService{
		userRepo:   userRepo,
		logService: logService,
		jwtSecret:  jwtSecret,
		jwtExpiry:  jwtExpiry,
	}
}

func (s *authService) buildAuthResponse(user *models.User) (models.AuthTokenResponse, error) {
	token, expiresAt, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return models.AuthTokenResponse{}, err
	}

	permissions := auth.PermissionsForRole(user.Role)
	expiresIn := int(time.Until(expiresAt).Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}

	return models.AuthTokenResponse{
		Token:       token,
		ExpiresAt:   expiresAt,
		ExpiresIn:   expiresIn,
		Role:        user.Role,
		Permissions: permissions,
		User:        models.ToUserResponse(user),
	}, nil
}

func (s *authService) Login(req models.LoginRequest) (models.AuthTokenResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.AuthTokenResponse{}, ErrInvalidCredentials
		}
		return models.AuthTokenResponse{}, err
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return models.AuthTokenResponse{}, ErrInvalidCredentials
	}

	resp, err := s.buildAuthResponse(user)
	if err != nil {
		return models.AuthTokenResponse{}, err
	}

	s.logService.LogAsync(user.ID, models.LogEventUserLogin, map[string]interface{}{
		"email": user.Email,
	})

	return resp, nil
}

func (s *authService) Register(req models.RegisterRequest) (models.AuthTokenResponse, error) {
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return models.AuthTokenResponse{}, ErrEmailTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.AuthTokenResponse{}, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return models.AuthTokenResponse{}, err
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     models.RoleUser,
	}

	if err := s.userRepo.Create(user); err != nil {
		return models.AuthTokenResponse{}, err
	}

	resp, err := s.buildAuthResponse(user)
	if err != nil {
		return models.AuthTokenResponse{}, err
	}

	s.logService.LogAsync(user.ID, models.LogEventUserRegister, map[string]interface{}{
		"email": user.Email,
		"name":  user.Name,
	})

	return resp, nil
}

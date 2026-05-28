package services

import (
	"errors"

	"user-management-api/models"
	"user-management-api/repositories"
	"user-management-api/utils"

	"gorm.io/gorm"
)

type UserService interface {
	Create(req models.CreateUserRequest, actorID uint) (models.UserResponse, error)
	GetByID(id uint) (models.UserResponse, error)
	Update(id uint, req models.UpdateUserRequest, actorID uint) (models.UserResponse, error)
	Delete(id uint, actorID uint) error
	List(page, perPage int, search string) ([]models.UserResponse, int64, error)
}

type userService struct {
	userRepo   repositories.UserRepository
	logService LogService
}

func NewUserService(userRepo repositories.UserRepository, logService LogService) UserService {
	return &userService{
		userRepo:   userRepo,
		logService: logService,
	}
}

func (s *userService) Create(req models.CreateUserRequest, actorID uint) (models.UserResponse, error) {
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return models.UserResponse{}, ErrEmailTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.UserResponse{}, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return models.UserResponse{}, err
	}

	role := models.RoleUser
	if req.Role != "" {
		role = req.Role
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return models.UserResponse{}, err
	}

	s.logService.LogAsync(actorID, models.LogEventUserCreated, map[string]interface{}{
		"target_user_id": user.ID,
		"email":          user.Email,
		"name":           user.Name,
	})

	return models.ToUserResponse(user), nil
}

func (s *userService) GetByID(id uint) (models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UserResponse{}, ErrUserNotFound
		}
		return models.UserResponse{}, err
	}
	return models.ToUserResponse(user), nil
}

func (s *userService) Update(id uint, req models.UpdateUserRequest, actorID uint) (models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UserResponse{}, ErrUserNotFound
		}
		return models.UserResponse{}, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" && req.Email != user.Email {
		if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
			return models.UserResponse{}, ErrEmailTaken
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UserResponse{}, err
		}
		user.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return models.UserResponse{}, err
		}
		user.Password = hashedPassword
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := s.userRepo.Update(user); err != nil {
		return models.UserResponse{}, err
	}

	s.logService.LogAsync(actorID, models.LogEventUserUpdated, map[string]interface{}{
		"target_user_id": user.ID,
		"email":          user.Email,
	})

	return models.ToUserResponse(user), nil
}

func (s *userService) Delete(id uint, actorID uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if err := s.userRepo.Delete(id); err != nil {
		return err
	}

	s.logService.LogAsync(actorID, models.LogEventUserDeleted, map[string]interface{}{
		"target_user_id": id,
		"email":          user.Email,
	})

	return nil
}

func (s *userService) List(page, perPage int, search string) ([]models.UserResponse, int64, error) {
	users, total, err := s.userRepo.List(page, perPage, search)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.UserResponse, len(users))
	for i, u := range users {
		responses[i] = models.ToUserResponse(&u)
	}

	return responses, total, nil
}

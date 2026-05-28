package services

import (
	"errors"

	"user-management-api/models"
	"user-management-api/repositories"

	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryService interface {
	Create(req models.CreateCategoryRequest, actorID uint) (*models.Category, error)
	GetByID(id uint) (*models.Category, error)
	Update(id uint, req models.UpdateCategoryRequest, actorID uint) (*models.Category, error)
	Delete(id uint, actorID uint) error
	List(page, perPage int, search string) ([]models.Category, int64, error)
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
	logService   LogService
}

func NewCategoryService(categoryRepo repositories.CategoryRepository, logService LogService) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
		logService:   logService,
	}
}

func (s *categoryService) Create(req models.CreateCategoryRequest, actorID uint) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	s.logService.LogAsync(actorID, models.LogEventCategoryCreate, map[string]interface{}{
		"category_id": category.ID,
		"name":        category.Name,
	})

	return category, nil
}

func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return category, nil
}

func (s *categoryService) Update(id uint, req models.UpdateCategoryRequest, actorID uint) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	s.logService.LogAsync(actorID, models.LogEventCategoryUpdate, map[string]interface{}{
		"category_id": category.ID,
		"name":        category.Name,
	})

	return category, nil
}

func (s *categoryService) Delete(id uint, actorID uint) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}

	if err := s.categoryRepo.Delete(id); err != nil {
		return err
	}

	s.logService.LogAsync(actorID, models.LogEventCategoryDelete, map[string]interface{}{
		"category_id": id,
		"name":        category.Name,
	})

	return nil
}

func (s *categoryService) List(page, perPage int, search string) ([]models.Category, int64, error) {
	return s.categoryRepo.List(page, perPage, search)
}

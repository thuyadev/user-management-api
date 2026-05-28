package services

import (
	"errors"

	"user-management-api/models"
	"user-management-api/repositories"

	"gorm.io/gorm"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductService interface {
	Create(req models.CreateProductRequest, actorID uint) (*models.Product, error)
	GetByID(id uint) (*models.Product, error)
	Update(id uint, req models.UpdateProductRequest, actorID uint) (*models.Product, error)
	Delete(id uint, actorID uint) error
	List(page, perPage int, search string, categoryID uint) ([]models.Product, int64, error)
}

type productService struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
	logService   LogService
}

func NewProductService(
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
	logService LogService,
) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		logService:   logService,
	}
}

func (s *productService) Create(req models.CreateProductRequest, actorID uint) (*models.Product, error) {
	if _, err := s.categoryRepo.FindByID(req.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	s.logService.LogAsync(actorID, models.LogEventProductCreate, map[string]interface{}{
		"product_id": product.ID,
		"name":       product.Name,
	})

	return s.productRepo.FindByID(product.ID)
}

func (s *productService) GetByID(id uint) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (s *productService) Update(id uint, req models.UpdateProductRequest, actorID uint) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if req.CategoryID > 0 {
		if _, err := s.categoryRepo.FindByID(req.CategoryID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrCategoryNotFound
			}
			return nil, err
		}
		product.CategoryID = req.CategoryID
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	s.logService.LogAsync(actorID, models.LogEventProductUpdate, map[string]interface{}{
		"product_id": product.ID,
		"name":       product.Name,
	})

	return s.productRepo.FindByID(product.ID)
}

func (s *productService) Delete(id uint, actorID uint) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	if err := s.productRepo.Delete(id); err != nil {
		return err
	}

	s.logService.LogAsync(actorID, models.LogEventProductDelete, map[string]interface{}{
		"product_id": id,
		"name":       product.Name,
	})

	return nil
}

func (s *productService) List(page, perPage int, search string, categoryID uint) ([]models.Product, int64, error) {
	return s.productRepo.List(page, perPage, search, categoryID)
}

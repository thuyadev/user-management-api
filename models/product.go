package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:255;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	Stock       int            `json:"stock" gorm:"not null;default:0"`
	CategoryID  uint           `json:"category_id" gorm:"not null;index"`
	Category    *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Product) TableName() string {
	return "products"
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=255"`
	Description string  `json:"description" binding:"omitempty,max=2000"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	CategoryID  uint    `json:"category_id" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name        string   `json:"name" binding:"omitempty,min=2,max=255"`
	Description string   `json:"description" binding:"omitempty,max=2000"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       *int     `json:"stock" binding:"omitempty,gte=0"`
	CategoryID  uint     `json:"category_id" binding:"omitempty,gt=0"`
}

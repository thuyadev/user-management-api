package seeders

import (
	"log"

	"user-management-api/models"

	"gorm.io/gorm"
)

type CategorySeeder struct{}

func NewCategorySeeder() *CategorySeeder {
	return &CategorySeeder{}
}

func (s *CategorySeeder) Run(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Category{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Println("CategorySeeder: categories already exist, skipping")
		return nil
	}

	categories := []models.Category{
		{Name: "Electronics", Description: "Electronic devices and gadgets"},
		{Name: "Clothing", Description: "Apparel and fashion items"},
		{Name: "Books", Description: "Books and publications"},
		{Name: "Home & Garden", Description: "Home improvement and garden supplies"},
	}

	if err := db.Create(&categories).Error; err != nil {
		return err
	}

	log.Printf("CategorySeeder: seeded %d categories", len(categories))
	return nil
}

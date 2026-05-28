package seeders

import (
	"log"

	"user-management-api/models"

	"gorm.io/gorm"
)

type ProductSeeder struct{}

func NewProductSeeder() *ProductSeeder {
	return &ProductSeeder{}
}

func (s *ProductSeeder) Run(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Product{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Println("ProductSeeder: products already exist, skipping")
		return nil
	}

	products := []models.Product{
		{Name: "Wireless Headphones", Description: "Noise-cancelling Bluetooth headphones", Price: 99.99, Stock: 50, CategoryID: 1},
		{Name: "Smartphone", Description: "Latest generation smartphone", Price: 699.99, Stock: 30, CategoryID: 1},
		{Name: "T-Shirt", Description: "Cotton casual t-shirt", Price: 19.99, Stock: 100, CategoryID: 2},
		{Name: "Jeans", Description: "Classic fit denim jeans", Price: 49.99, Stock: 75, CategoryID: 2},
		{Name: "Go Programming Book", Description: "Learn Go programming language", Price: 39.99, Stock: 40, CategoryID: 3},
		{Name: "Garden Hose", Description: "50ft expandable garden hose", Price: 29.99, Stock: 60, CategoryID: 4},
	}

	if err := db.Create(&products).Error; err != nil {
		return err
	}

	log.Printf("ProductSeeder: seeded %d products", len(products))
	return nil
}

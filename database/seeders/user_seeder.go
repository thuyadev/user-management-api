package seeders

import (
	"log"

	"user-management-api/models"
	"user-management-api/utils"

	"gorm.io/gorm"
)

type UserSeeder struct {
	cfg *utils.Config
}

func NewUserSeeder(cfg *utils.Config) *UserSeeder {
	return &UserSeeder{cfg: cfg}
}

func (s *UserSeeder) Run(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Println("UserSeeder: users already exist, skipping")
		return nil
	}

	hashedPassword, err := utils.HashPassword(s.cfg.AdminPassword)
	if err != nil {
		return err
	}

	users := []models.User{
		{
			Name:     s.cfg.AdminName,
			Email:    s.cfg.AdminEmail,
			Password: hashedPassword,
			Role:     models.RoleAdmin,
		},
		{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: mustHash("password123"),
			Role:     models.RoleUser,
		},
		{
			Name:     "Jane Smith",
			Email:    "jane@example.com",
			Password: mustHash("password123"),
			Role:     models.RoleUser,
		},
	}

	if err := db.Create(&users).Error; err != nil {
		return err
	}

	log.Printf("UserSeeder: seeded %d users", len(users))
	return nil
}

func mustHash(password string) string {
	hash, err := utils.HashPassword(password)
	if err != nil {
		panic(err)
	}
	return hash
}

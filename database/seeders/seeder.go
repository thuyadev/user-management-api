package seeders

import (
	"log"

	"user-management-api/utils"

	"gorm.io/gorm"
)

type Seeder interface {
	Run(db *gorm.DB) error
}

type DatabaseSeeder struct {
	seeders []Seeder
}

func NewDatabaseSeeder(seeders ...Seeder) *DatabaseSeeder {
	return &DatabaseSeeder{seeders: seeders}
}

func (ds *DatabaseSeeder) Run(db *gorm.DB) error {
	for _, seeder := range ds.seeders {
		if err := seeder.Run(db); err != nil {
			return err
		}
	}
	log.Println("All seeders completed successfully")
	return nil
}

func RunAll(db *gorm.DB, cfg *utils.Config) error {
	seeder := NewDatabaseSeeder(
		NewUserSeeder(cfg),
		NewCategorySeeder(),
		NewProductSeeder(),
	)
	return seeder.Run(db)
}

package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Schema struct {
	gorm.Model
	Id string `json:"id"`
	Version int `json:"version"`
	LatestVersion int `json:"latestVersion"`
}

func InitDb(path string) (*gorm.DB, error) {
	// Open db
	db, dbErr := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if dbErr != nil {
		return nil, dbErr
	}

	// Migrate db
	migrateErr := db.AutoMigrate(&Schema{})
	if migrateErr != nil {
		return nil, migrateErr
	}

	// Return db
	return db, nil
}
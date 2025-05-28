package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Schema struct {
	gorm.Model
	Id            int    `json:"id" gorm:"primaryKey"`
	Urn           string `json:"urn"`
	Version       int    `json:"version"`
	LatestVersion int    `json:"latestVersion"`
}

func InitDatabase(path string) (*gorm.DB, error) {
	// Open db
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, err
	}

	// Migrate db
	err = db.AutoMigrate(&Schema{})

	if err != nil {
		return nil, err
	}

	// Return db
	return db, nil
}

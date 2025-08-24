package database

import (
	"embed"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Apps struct {
	gorm.Model
	Urn           string `json:"urn"`
	Version       int    `json:"version"`
	LatestVersion int    `json:"latestVersion"`
}

//go:embed migrations/*.sql
var migrations embed.FS

func migrate(db *gorm.DB) error {
	files, err := migrations.ReadDir("migrations")

	if err != nil {
		return err
	}

	// Execute each migration file
	for _, file := range files {
		content, err := migrations.ReadFile("migrations/" + file.Name())

		if err != nil {
			return err
		}

		err = db.Exec(string(content)).Error

		if err != nil {
			return err
		}
	}

	return nil
}

func InitDatabase(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, err
	}

	err = migrate(db)

	if err != nil {
		return nil, err
	}

	return db, nil
}

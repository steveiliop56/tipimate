package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Apps struct {
	gorm.Model
	Urn           string
	Version       int
	LatestVersion int
}

type AppsOld struct {
	gorm.Model
	Id            string
	Version       int
	LatestVersion int
}

func InitDatabase(path string) (*gorm.DB, error) {
	// Open db
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, err
	}

	// Rename old table if it exists
	if db.Migrator().HasTable("schemas") {
		err = db.Migrator().RenameTable("schemas", "apps_old")
		if err != nil {
			return nil, err
		}
	}

	// Migrate db
	err = db.AutoMigrate(&Apps{})

	if err != nil {
		return nil, err
	}

	// Migrate old data
	var oldRecords []AppsOld
	db.Table("apps_old").Find(&oldRecords)

	for _, record := range oldRecords {
		// Save new record
		res := db.Create(&Apps{
			Urn:           record.Id + ":migrated",
			Version:       record.Version,
			LatestVersion: record.LatestVersion,
		})
		if res.Error != nil {
			return nil, res.Error
		}
	}

	// Drop old table
	if db.Migrator().HasTable("apps_old") {
		err = db.Migrator().DropTable("apps_old")
		if err != nil {
			return nil, err
		}
	}

	// Return db
	return db, nil
}

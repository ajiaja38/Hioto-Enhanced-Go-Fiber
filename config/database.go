package config

import (
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func DBConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(DB_PATH.GetValue()+"?_loc=Asia%2FJakarta"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	log.Info("Successfully connected to SQLite database 🗃️")

	return db, nil
}

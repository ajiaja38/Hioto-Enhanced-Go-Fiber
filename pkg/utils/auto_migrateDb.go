package utils

import (
	"go/hioto/pkg/model"

	"gorm.io/gorm"
)

func AutoMigrateDb(db *gorm.DB) {
	db.AutoMigrate(&model.Registration{})
	db.AutoMigrate(&model.RuleDevice{})
	db.AutoMigrate(&model.Log{})
	db.AutoMigrate(&model.LogAktuator{})
	db.AutoMigrate(&model.MonitoringHistory{})
}

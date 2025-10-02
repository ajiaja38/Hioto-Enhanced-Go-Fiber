package router

import (
	"go/hioto/pkg/handler/res"
	"go/hioto/pkg/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ControlDeviceRouter(router fiber.Router, db *gorm.DB, controlDeviceService *service.ControlDeviceService) {
	controlDeviceHandler := res.NewControlDeviceHandler(controlDeviceService)

	router.Put("/device/control", controlDeviceHandler.ControlDeviceHandler)
}

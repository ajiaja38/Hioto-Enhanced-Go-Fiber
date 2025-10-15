package router

import (
	"go/hioto/pkg/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Router(
	router fiber.Router,
	db *gorm.DB,
	controlDeviceService *service.ControlDeviceService,
	deviceService *service.DeviceService,
	rulesService *service.RuleService,
) {
	ControlDeviceRouter(router, db, controlDeviceService)
	DeviceRouter(router, db, deviceService)
	RulesRouter(router, db, rulesService)
}

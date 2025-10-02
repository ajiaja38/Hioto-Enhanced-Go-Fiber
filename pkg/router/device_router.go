package router

import (
	"go/hioto/pkg/handler/res"
	"go/hioto/pkg/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DeviceRouter(router fiber.Router, db *gorm.DB, deviceService *service.DeviceService) {
	deviceHandler := res.NewDeviceHandler(deviceService)

	router.Post("/device", deviceHandler.RegisterDevice)
	router.Get("/devices", deviceHandler.GetAllDeviceHandler)
	router.Get("/device/:guid", deviceHandler.GetDeviceByGuidHandler)
	router.Put("/device", deviceHandler.UpdateDeviceByGuidHandler)
	router.Delete("/device/:guid", deviceHandler.DeleteDeviceByGuidHandler)
}

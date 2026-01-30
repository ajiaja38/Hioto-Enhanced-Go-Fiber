package router

import (
	"go/hioto/pkg/handler/res"
	"go/hioto/pkg/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FloorRouter(router fiber.Router, db *gorm.DB, floorService *service.FloorService) {
	floorHandler := res.NewFloorHandler(floorService)

	router.Post("/floor", floorHandler.CreateFloor)
	router.Get("/floors", floorHandler.GetAllFloors)
	router.Get("/floor/:id", floorHandler.GetFloorByID)
	router.Put("/floor/:id", floorHandler.UpdateFloor)
	router.Delete("/floor/:id", floorHandler.DeleteFloor)
}

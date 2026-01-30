package router

import (
	"go/hioto/pkg/handler/res"
	"go/hioto/pkg/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RoomRouter(router fiber.Router, db *gorm.DB, roomService *service.RoomService) {
	roomHandler := res.NewRoomHandler(roomService)

	router.Post("/room", roomHandler.CreateRoom)
	router.Get("/rooms", roomHandler.GetAllRooms)
	router.Get("/rooms/floor/:floor_id", roomHandler.GetRoomsByFloorID)
	router.Get("/room/:id", roomHandler.GetRoomByID)
	router.Put("/room/:id", roomHandler.UpdateRoom)
	router.Delete("/room/:id", roomHandler.DeleteRoom)
}

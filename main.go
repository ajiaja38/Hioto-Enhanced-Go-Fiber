package main

import (
	"context"
	"go/hioto/config"
	"go/hioto/pkg/handler/consumer"
	"go/hioto/pkg/handler/err"
	"go/hioto/pkg/router"
	"go/hioto/pkg/service"
	"go/hioto/pkg/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config.Load()

	db, errDb := config.DBConnection()

	if errDb != nil {
		log.Fatalf("Failed to connect to database %s: %v", config.DB_PATH.GetValue(), errDb)
	}

	utils.AutoMigrateDb(db)

	config.CreateRmqInstance()
	defer config.CloseRabbitMQ()

	config.CreateMqttInstance()
	defer config.CloseAllMqttInstances()

	// Start Workers
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// define instance services
	controlDeviceService := service.NewControlDeviceService(db)
	deviceService := service.NewDeviceService(db)
	ruleService := service.NewRuleService(db)
	floorService := service.NewFloorService(db)
	roomService := service.NewRoomService(db)

	// Start Consumer
	consumerHandler := consumer.NewConsumerHandler(ruleService, deviceService, controlDeviceService)
	consumerRouter := router.NewConsumerMessageBroker(ctx, consumerHandler)
	consumerRouter.StartConsumer()

	log.Info("Hello From Hioto Worker 💡")

	// REST API FIBER
	app := fiber.New(fiber.Config{
		ErrorHandler: err.ErrorHandler,
	})

	port := config.PORT.GetValue()

	if port == "" {
		port = "8080"
	}

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
	app.Use(recover.New())

	route := app.Group("/api")
	route.Get("/", func(c *fiber.Ctx) error {
		return utils.SuccessResponse[*struct{}](c, fiber.StatusOK, "Hi From API Local Hioto's 💡", nil)
	})
	route.Get("/metrics", monitor.New(monitor.Config{Title: "Hioto Metrics Pages"}))

	// REST API Router Group
	router.Router(route, db, controlDeviceService, deviceService, ruleService, floorService, roomService)

	log.Infof("API server is running on http://localhost:%s/api 💡", port)

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("http server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Warn("Shutting down gracefully...💡")

	if err := app.Shutdown(); err != nil {
		log.Errorf("Error shutting down Fiber: %v", err)
	}

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

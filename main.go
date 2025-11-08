package main

import (
	"context"
	"go/hioto/config"
	"go/hioto/pkg/handler/consumer"
	"go/hioto/pkg/handler/err"
	"go/hioto/pkg/router"
	"go/hioto/pkg/service"
	"go/hioto/pkg/utils"
	"go/hioto/pkg/utils/cron"
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
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Error("Error loading .env file")
	}
}

func main() {
	db, errDb := config.DBConnection()

	if errDb != nil {
		log.Fatalf("Failed to connect to database: %v", errDb)
	}

	utils.AutoMigrateDb(db)

	// Start Worker
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// define instance service
	controlDeviceService := service.NewControlDeviceService(db)
	deviceService := service.NewDeviceService(db)
	ruleService := service.NewRuleService(db)
	logService := service.NewLogService(db)

	// Start Consumer
	consumerHandler := consumer.NewConsumerHandler(ruleService, deviceService, controlDeviceService)
	consumerRouter := router.NewConsumerMessageBroker(consumerHandler, ctx)
	consumerRouter.StartConsumer()

	cron.LoggerCrobJob(*logService)

	log.Info("Hello From Worker Hioto💡")

	// REST API FIBER
	app := fiber.New(fiber.Config{
		ErrorHandler: err.ErrorHandler,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
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
		return utils.SuccessResponse[*struct{}](c, fiber.StatusOK, "Hello From API Local Gokils Hioto 💡", nil)
	})
	route.Get("/metrics", monitor.New(monitor.Config{Title: "Fiber Metrics Page"}))

	// REST API Router Group
	router.Router(route, db, controlDeviceService, deviceService, ruleService)

	log.Infof("Starting server on http://localhost:%s/api 💡", port)

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Warn("Shutting down gracefully...")

	if err := app.Shutdown(); err != nil {
		log.Errorf("Error shutting down Fiber: %v", err)
	}

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

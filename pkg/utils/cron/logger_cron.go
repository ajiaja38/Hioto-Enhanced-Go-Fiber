package cron

import (
	"go/hioto/pkg/service"

	"github.com/gofiber/fiber/v2/log"
	"github.com/robfig/cron/v3"
)

func LoggerCrobJob(logService service.LogService) {
	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	c.AddFunc("@every 10m", func() {
		logService.GetAllLogs()
		logService.GetAllLogAktuators()
		logService.GetAllMonitoringHistory()
	})

	log.Info("Starting cron job...")

	c.Start()
}

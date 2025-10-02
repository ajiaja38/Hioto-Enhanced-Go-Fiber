package router

import (
	"go/hioto/pkg/handler/res"
	"go/hioto/pkg/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RulesRouter(router fiber.Router, db *gorm.DB, rulesService *service.RuleService) {
	rulesHandler := res.NewRulesHandler(rulesService)

	router.Post("/rule", rulesHandler.CreateRulesHandler)
	router.Get("/rule/:guidDevice", rulesHandler.GetRulesByGuidHandler)
	router.Delete("/rule/:guidSensor", rulesHandler.DeleteRulesByGuidSensorHandler)
}

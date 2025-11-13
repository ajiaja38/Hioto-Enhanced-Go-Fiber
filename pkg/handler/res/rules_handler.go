package res

import (
	"go/hioto/pkg/dto"
	"go/hioto/pkg/service"
	"go/hioto/pkg/utils"
	"go/hioto/pkg/utils/validators"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type RulesHandler struct {
	rulesService *service.RuleService
	validator    *validator.Validate
}

func NewRulesHandler(rulesService *service.RuleService) *RulesHandler {
	validate := validator.New()
	validators.RegisterCustomValidators(validate)

	return &RulesHandler{
		rulesService: rulesService,
		validator:    validate,
	}
}

func (h *RulesHandler) CreateRulesHandler(c *fiber.Ctx) error {
	var createRuleDto dto.CreateRuleDto

	if err := utils.ValidateRequestBody(c, h.validator, &createRuleDto); err != nil {
		return err
	}

	responseRules, err := h.rulesService.CreateRules(&createRuleDto)

	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Success create rules", responseRules)
}

func (h *RulesHandler) GetRulesByGuidHandler(c *fiber.Ctx) error {
	guid := c.Params("guidDevice")

	rules, err := h.rulesService.GetRulesByGuid(guid)

	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Success get rules by guid", rules)
}

func (h *RulesHandler) DeleteRulesByGuidSensorHandler(c *fiber.Ctx) error {
	guid := c.Params("guidSensor")

	if err := h.rulesService.DeleteRulesByGuidSensor(guid); err != nil {
		return err
	}

	return utils.SuccessResponse[any](c, fiber.StatusOK, "Success delete rules by guid sensor", nil)
}

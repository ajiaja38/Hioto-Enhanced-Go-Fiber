package utils

import (
	"go/hioto/pkg/model"

	"github.com/gofiber/fiber/v2"
)

func SuccessResponse[T any](c *fiber.Ctx, code int, message string, data T) error {
	return c.Status(code).JSON(model.ResponseEntity[T]{
		Code:    code,
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func SuccessResponsePaginate[T any](c *fiber.Ctx, code int, message string, data T, meta *model.MetaPagination) error {
	return c.Status(code).JSON(model.ResponseEntity[T]{
		Code:    code,
		Status:  true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

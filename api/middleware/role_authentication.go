package middleware

import (
	"shopping-site/pkg/loggers"
	"shopping-site/utils/constants"
	"shopping-site/utils/dto"

	"github.com/gofiber/fiber/v2"
)

func AdminRoleAuthentication(ctx *fiber.Ctx) error {
	role := ctx.Locals("role")
	if role == "" {
		loggers.ErrorLog.Println("role authentication is empty")
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ResponseJson{
			Error: "insuffisient permission",
		})
	}

	if role != constants.AdminRole {
		loggers.ErrorLog.Println("insufficient permission")
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.ResponseJson{
			Error: "insuffisient permission",
		})
	}
	return ctx.Next()
}

func MerchantRoleAuthentication(ctx *fiber.Ctx) error {
	role := ctx.Locals("role")
	if role == "" {
		loggers.ErrorLog.Println("role authentication is empty")
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ResponseJson{
			Error: "insuffisient permission",
		})
	}

	if role != constants.MerchantRole {
		loggers.ErrorLog.Println("insufficient permission")
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.ResponseJson{
			Error: "insuffisient permission",
		})
	}
	return ctx.Next()
}

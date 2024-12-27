package handlers

import (
	"shopping-site/api/services"
	"shopping-site/pkg/loggers"
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	services.IAdminService
}

func (service *AdminHandler) AddCategoreyHandler(ctx *fiber.Ctx) error {
	var category models.Categories

	if err := ctx.BodyParser(&category); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	errResponse := service.IAdminService.AddCategoreyService(&category)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{
		Message: "Category created successfully",
		Data:    category,
	})
}

func (service *AdminHandler) AddBrandHandler(ctx *fiber.Ctx) error {
	var brand models.Brands

	if err := ctx.BodyParser(&brand); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	errResponse := service.IAdminService.AddBrandService(&brand)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{
		Message: "Brand created successfully",
		Data:    brand,
	})
}

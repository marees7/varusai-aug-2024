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

// create new category
//
//	@Summary		Create category
//	@Description	Create a new category
//	@ID				create_category
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			Category	details		body	models.Categories	true	"Enter category details"
//	@Success		200			{object}	dto.ResponseJson
//	@Failure		400			{object}	dto.ResponseJson
//	@Failure		409			{object}	dto.ResponseJson
//	@Router			/admin/category [post]
func (service *AdminHandler) CreateCategorey(ctx *fiber.Ctx) error {
	var category models.Categories

	if err := ctx.BodyParser(&category); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	// call create category service
	errResponse := service.IAdminService.CreateCategorey(&category)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{
		Message: "category created successfully",
		Data:    category,
	})
}

// create new brand
//
//	@Summary		Create brand
//	@Description	Create a new brand
//	@ID				create_brand
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			Login	body		models.Brands	true	"Enter brand details"
//	@Success		200		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		409		{object}	dto.ResponseJson
//	@Router			/admin/brand [post]
func (service *AdminHandler) CreateBrand(ctx *fiber.Ctx) error {
	var brand models.Brands

	if err := ctx.BodyParser(&brand); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	// call create brand service
	errResponse := service.IAdminService.CreateBrand(&brand)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{
		Message: "brand created successfully",
		Data:    brand,
	})
}

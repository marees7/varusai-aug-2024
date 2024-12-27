package handlers

import (
	"fmt"
	"shopping-site/api/services"
	"shopping-site/pkg/loggers"
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MerchantHandler struct {
	services.IMerchantService
}

func (service *MerchantHandler) AddProductHandler(ctx *fiber.Ctx) error {
	var product models.Products
	userIdCtx := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&product); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	errResponse := service.IMerchantService.AddProductService(userIdCtx, &product)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{
		Message: "product created successfully",
		Data:    product,
	})
}

func (service *MerchantHandler) RemoveProductHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	errResponse := service.IMerchantService.RemoveProductService(id)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "product Deleted successfully",
	})
}

func (service *MerchantHandler) UpdateProductHandler(ctx *fiber.Ctx) error {
	var product models.Products
	userIdCtx := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&product); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	errResponse := service.IMerchantService.UpdateProductService(userIdCtx, &product)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{
		Message: "product updated successfully",
		Data:    map[string]interface{}{"product_id": product.ProductId},
	})
}

func (service *MerchantHandler) UpdateMerchantHandler(ctx *fiber.Ctx) error {
	var user models.Users
	userIdCtx := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&user); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	errResponse := service.IMerchantService.UpdateMerchantService(userIdCtx, &user)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{
		Message: "user details updated successfully",
		Data:    map[string]interface{}{"user_id": user.UserId},
	})
}

func (service *MerchantHandler) UpdateOrderStatusHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	orderStatus := ctx.Query("order_status")
	fmt.Println("**", orderStatus)
	userIdCtx := ctx.Locals("user_id").(uuid.UUID)

	errResponse := service.IMerchantService.UpdateOrderStatusService(userIdCtx, id, orderStatus)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{
		Message: "order status updated successfully",
		Data:    map[string]interface{}{"order_id": id},
	})
}

func (service *MerchantHandler) GetProductsHandler(ctx *fiber.Ctx) error {
	userIdCtx := ctx.Locals("user_id").(uuid.UUID)

	filters := ctx.Queries()

	products, errResponse := service.IMerchantService.GetProductsService(filters, userIdCtx)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Data: products,
	})
}

func (service *MerchantHandler) GetProductHandler(ctx *fiber.Ctx) error {
	userIdCtx := ctx.Locals("user_id").(uuid.UUID)
	id := ctx.Params("id")

	product, errResponse := service.IMerchantService.GetProductService(userIdCtx, id)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{
		Data: product,
	})
}

func (service *MerchantHandler) GetOrdersHandler(ctx *fiber.Ctx) error {
	userIdCtx := ctx.Locals("user_id").(uuid.UUID)

	orders, errResponse := service.IMerchantService.GetOrdersService(userIdCtx)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{
		Data: orders,
	})
}

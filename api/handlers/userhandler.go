package handlers

import (
	"shopping-site/api/services"
	"shopping-site/pkg/loggers"
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	services.IUserService
}

func (service *UserHandler) PlaceOrderHandler(ctx *fiber.Ctx) error {
	var order models.Orders
	userIdCtx := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&order); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	orderDetails, errResponse := service.IUserService.PlaceOrderService(userIdCtx, order)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{
		Message: "order placed successfully",
		Data:    orderDetails,
	})
}

func (service *UserHandler) CancelOrderHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userIdCtx := ctx.Locals("user_id").(uuid.UUID)

	errResponse := service.IUserService.CancelOrderService(userIdCtx, id)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "order cancelled successfully",
	})
}

func (service *UserHandler) UpdateUserHandler(ctx *fiber.Ctx) error {
	var user models.Users
	userIdCtx := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&user); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	errResponse := service.IUserService.UpdateUserService(userIdCtx, &user)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "user details updated successfully",
		Data:    map[string]interface{}{"user_id": user.UserId},
	})
}

func (service *UserHandler) GetOrdersHandler(ctx *fiber.Ctx) error {
	userIdCtx := ctx.Locals("user_id").(uuid.UUID)

	orders, errResponse := service.IUserService.GetOrdersService(userIdCtx)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Data: orders,
	})
}

func (service *UserHandler) GetProductsHandler(ctx *fiber.Ctx) error {
	userIdCtx := ctx.Locals("user_id").(uuid.UUID)

	filters := ctx.Queries()

	products, errResponse := service.IUserService.GetProductsService(filters, userIdCtx)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Data: products,
	})
}

func (service *UserHandler) GetProductHandler(ctx *fiber.Ctx) error {
	userIdCtx := ctx.Locals("user_id").(uuid.UUID)
	id := ctx.Params("id")

	product, errResponse := service.IUserService.GetProductService(userIdCtx, id)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Data: product,
	})
}

func (service *UserHandler) FilterProductsHandler(ctx *fiber.Ctx) error {
	filters := ctx.Queries()

	products, errResponse := service.IUserService.FilterProductsService(filters)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Data: products,
	})
}

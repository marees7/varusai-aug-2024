package handlers

import (
	"shopping-site/api/services"
	"shopping-site/pkg/loggers"
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"
	"shopping-site/utils/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	services.IUserService
}

func (service *UserHandler) CreateOrder(ctx *fiber.Ctx) error {
	var order models.Orders
	userId := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&order); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	orderDetails, errResponse := service.IUserService.PlaceOrder(userId, order)
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

func (service *UserHandler) GetProducts(ctx *fiber.Ctx) error {
	var (
		price  float64
		rating float32
	)
	filters := make(map[string]interface{})
	keywords := ctx.Queries()

	limitstr := keywords["limit"]
	offsetstr := keywords["offset"]

	limit, offset, err := helper.Pagination(limitstr, offsetstr)
	if err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	if keywords["price"] != "" {
		price, err = helper.ToFloat(keywords["price"])
		if err != nil {
			loggers.WarnLog.Println(err)
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
				Error: err.Error(),
			})
		}
	} else if keywords["rating"] != "" {
		result, err := helper.ToFloat(keywords["rating"])
		if err != nil {
			loggers.WarnLog.Println(err)
			return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
				Error: err.Error(),
			})
		}

		rating = float32(result)
	}

	filters["limit"] = limit
	filters["offset"] = offset
	filters["category_name"] = keywords["category_name"]
	filters["brand_name"] = keywords["brand_name"]
	filters["price"] = price
	filters["rating"] = rating

	products, count, errResponse := service.IUserService.GetProducts(filters)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Data:         products,
		TotalRecords: count,
		Limit:        limit,
		Offset:       offset,
	})
}

func (service *UserHandler) GetProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	product, errResponse := service.IUserService.GetProduct(id)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Data: product,
	})
}

func (service *UserHandler) GetOrders(ctx *fiber.Ctx) error {
	filters := make(map[string]interface{})
	userId := ctx.Locals("user_id").(uuid.UUID)

	keywords := ctx.Queries()

	limitstr := keywords["limit"]
	offsetstr := keywords["offset"]

	limit, offset, err := helper.Pagination(limitstr, offsetstr)
	if err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	filters["limit"] = limit
	filters["offset"] = offset
	filters["from_date"] = keywords["from_date"]
	filters["to_date"] = keywords["to_date"]

	orders, count, errResponse := service.IUserService.GetOrders(userId, filters)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Data:         orders,
		TotalRecords: count,
		Limit:        limit,
		Offset:       offset,
	})
}

func (service *UserHandler) GetOrder(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(uuid.UUID)
	id := ctx.Params("id")

	order, errResponse := service.IUserService.GetOrder(userId, id)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Data: order,
	})
}

func (service *UserHandler) UpdateOrder(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userId := ctx.Locals("user_id").(uuid.UUID)

	errResponse := service.IUserService.UpdateOrder(userId, id)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "order cancelled successfully",
	})
}

func (service *UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	var user models.Users
	userId := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&user); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	errResponse := service.IUserService.UpdateUser(userId, &user)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "user details updated successfully",
		Data:    map[string]interface{}{"user_id": user.UserID},
	})
}

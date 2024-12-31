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

type MerchantHandler struct {
	services.IMerchantService
}

func (service *MerchantHandler) CreateProduct(ctx *fiber.Ctx) error {
	var product models.Products
	userId := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&product); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	errResponse := service.IMerchantService.CreateProduct(userId, &product)
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

func (service *MerchantHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	errResponse := service.IMerchantService.DeleteProduct(id)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "product deleted successfully",
	})
}

func (service *MerchantHandler) UpdateProduct(ctx *fiber.Ctx) error {
	var product models.Products
	userId := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&product); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	errResponse := service.IMerchantService.UpdateProduct(userId, &product)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "product updated successfully",
		Data:    map[string]interface{}{"product_id": product.ProductID},
	})
}

func (service *MerchantHandler) UpdateMerchant(ctx *fiber.Ctx) error {
	var user models.Users
	userId := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&user); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	errResponse := service.IMerchantService.UpdateMerchant(userId, &user)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "merchant details updated successfully",
		Data:    map[string]interface{}{"user_id": user.UserID},
	})
}

func (service *MerchantHandler) UpdateOrderStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	orderStatus := ctx.Query("order_status")
	userId := ctx.Locals("user_id").(uuid.UUID)

	errResponse := service.IMerchantService.UpdateOrderStatus(userId, id, orderStatus)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "order_item status updated successfully",
		Data:    map[string]interface{}{"order_item_id": id},
	})
}

func (service *MerchantHandler) GetProducts(ctx *fiber.Ctx) error {
	var (
		price  float64
		rating float32
	)
	filters := make(map[string]interface{})

	userId := ctx.Locals("user_id").(uuid.UUID)
	if userId == uuid.Nil {
		loggers.WarnLog.Println("user_id is empty")
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.ResponseJson{
			Error: "User_id is empty"})
	}

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

	products, count, errResponse := service.IMerchantService.GetProducts(filters, userId)
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

func (service *MerchantHandler) GetProduct(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(uuid.UUID)
	id := ctx.Params("id")

	product, errResponse := service.IMerchantService.GetProduct(userId, id)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Data: product,
	})
}

func (service *MerchantHandler) GetOrders(ctx *fiber.Ctx) error {
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

	orders, count, errResponse := service.IMerchantService.GetOrders(userId, filters)
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

func (service *MerchantHandler) GetOrder(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(uuid.UUID)
	id := ctx.Params("id")

	product, errResponse := service.IMerchantService.GetOrder(userId, id)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Data: product,
	})
}

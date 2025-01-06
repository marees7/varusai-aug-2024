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

// create new order
//
//	@Summary		Create order
//	@Description	Create a new order
//	@ID				create_order
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			CreateOrder	body		models.Orders	true	"Enter order details"
//	@Success		201			{object}	dto.ResponseJson
//	@Failure		400			{object}	dto.ResponseJson
//	@Failure		401			{object}	dto.ResponseJson
//	@Failure		403			{object}	dto.ResponseJson
//	@Router			/user/order [post]
func (service *UserHandler) CreateOrder(ctx *fiber.Ctx) error {
	var order models.Orders
	userId := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&order); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	// call create order service
	orderDetails, errResponse := service.IUserService.CreateOrder(userId, order)
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

// get products based on the filters
//
//	@Summary		Get products
//	@Description	Get all products
//	@ID				get_products
//	@Tags			Common
//	@Produce		json
//	@Security		JWT
//	@Param			limit			query		string	false	"Enter limit"
//	@Param			offset			query		string	false	"Enter offset"
//	@Param			category_name	query		string	false	"Enter category_name"
//	@Param			brand_name		query		string	false	"Enter brand_name"
//	@Param			price			query		string	false	"Enter price"
//	@Param			rating			query		string	false	"Enter rating"
//	@Success		200				{object}	dto.ResponseJson
//	@Failure		400				{object}	dto.ResponseJson
//	@Failure		401				{object}	dto.ResponseJson
//	@Failure		403				{object}	dto.ResponseJson
//	@Router			/common/product [get]
func (service *UserHandler) GetProducts(ctx *fiber.Ctx) error {
	var (
		price  float64
		rating float32
	)
	filters := make(map[string]interface{})

	keywords := ctx.Queries()

	limitstr := keywords["limit"]
	offsetstr := keywords["offset"]

	// pagination
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

	// call get products service
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

// get a single product by id
//
//	@Summary		Get product
//	@Description	Get product
//	@ID				get_product
//	@Tags			Common
//	@Produce		json
//	@Security		JWT
//	@Param			product_id	path		string	true	"Enter product id"
//	@Success		200			{object}	dto.ResponseJson
//	@Failure		400			{object}	dto.ResponseJson
//	@Failure		404			{object}	dto.ResponseJson
//	@Failure		401			{object}	dto.ResponseJson
//	@Failure		403			{object}	dto.ResponseJson
//	@Router			/common/product/{product_id} [get]
func (service *UserHandler) GetProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	// call get product service
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

// get the user's orders based on provided filters
//
//	@Summary		Get orders
//	@Description	Get orders of the user
//	@ID				get_orders_user
//	@Tags			user
//	@Produce		json
//	@Security		JWT
//	@Param			limit		query		string	false	"Enter limit"
//	@Param			offset		query		string	false	"Enter offset"
//	@Param			from_date	query		string	false	"Enter from_date"
//	@Param			to_date		query		string	false	"Enter to_date"
//	@Success		200			{object}	dto.ResponseJson
//	@Failure		400			{object}	dto.ResponseJson
//	@Failure		401			{object}	dto.ResponseJson
//	@Failure		403			{object}	dto.ResponseJson
//	@Router			/user/order [get]
func (service *UserHandler) GetOrders(ctx *fiber.Ctx) error {
	filters := make(map[string]interface{})
	userId := ctx.Locals("user_id").(uuid.UUID)

	keywords := ctx.Queries()

	limitstr := keywords["limit"]
	offsetstr := keywords["offset"]

	// pagination
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

	// call get orders repository
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

// get user's order
//
//	@Summary		Get order
//	@Description	Get order of the user
//	@ID				get_order_user
//	@Tags			user
//	@Produce		json
//	@Security		JWT
//	@Param			order_id	path		string	true	"Enter order id"
//	@Success		200			{object}	dto.ResponseJson
//	@Failure		400			{object}	dto.ResponseJson
//	@Failure		404			{object}	dto.ResponseJson
//	@Failure		401			{object}	dto.ResponseJson
//	@Failure		403			{object}	dto.ResponseJson
//	@Router			/user/order/{order_id} [get]
func (service *UserHandler) GetOrder(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(uuid.UUID)
	id := ctx.Params("id")

	// call get order service
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

// cancel order of user
//
//	@Summary		Update order
//	@Description	Update order status for the merchant
//	@ID				update_order_user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			order_id	path		string			true	"Enter order id"
//	@Param			OrderStatus	body		models.Orders	true	"Enter order status"
//	@Success		200			{object}	dto.ResponseJson
//	@Failure		304			{object}	dto.ResponseJson
//	@Failure		400			{object}	dto.ResponseJson
//	@Failure		404			{object}	dto.ResponseJson
//	@Failure		401			{object}	dto.ResponseJson
//	@Failure		403			{object}	dto.ResponseJson
//	@Router			/user/order/{order_id} [patch]
func (service *UserHandler) UpdateOrder(ctx *fiber.Ctx) error {
	var order models.Orders

	id := ctx.Params("id")
	userId := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&order); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	// call update order service
	errResponse := service.IUserService.UpdateOrder(userId, id, &order)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "order cancelled successfully",
	})
}

// update user details
//
//	@Summary		Update user
//	@Description	Update details of user
//	@ID				update_details_user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			UserDetails	body		models.Users	true	"Enter user details"
//	@Success		200			{object}	dto.ResponseJson
//	@Success		304			{object}	dto.ResponseJson
//	@Failure		400			{object}	dto.ResponseJson
//	@Failure		404			{object}	dto.ResponseJson
//	@Failure		401			{object}	dto.ResponseJson
//	@Failure		403			{object}	dto.ResponseJson
//	@Router			/user [patch]
func (service *UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	var user models.Users
	userId := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&user); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	// call update user service
	errResponse := service.IUserService.UpdateUser(userId, &user)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "user details updated successfully",
	})
}

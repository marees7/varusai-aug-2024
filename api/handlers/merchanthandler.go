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

// create product
//
//	@Summary		Create product
//	@Description	Create a new product
//	@ID				create_product
//	@Tags			Merchant
//	@Accept			json
//	@Produce		json
//	@Security 	    JWT
//	@Param			CreateProduct body models.Brands true "Enter product details"
//	@Success		201		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		409		{object}	dto.ResponseJson
//	@Failure		404		{object}	dto.ResponseJson
//	@Failure		401		{object}	dto.ResponseJson
//	@Failure		403		{object}	dto.ResponseJson
//	@Router			/merchant/product [post]
func (service *MerchantHandler) CreateProduct(ctx *fiber.Ctx) error {
	var product models.Products
	userId := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&product); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	// call create product service
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

// get the avilable categories
//
//	@Summary		Get categories
//	@Description	Get all categories
//	@ID				get_categories
//	@Tags			Merchant
//	@Produce		json
//	@Security 	    JWT
//	@Success		200		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		401		{object}	dto.ResponseJson
//	@Failure		403		{object}	dto.ResponseJson
//	@Router			/merchant/category [get]
func (service *MerchantHandler) GetCategories(ctx *fiber.Ctx) error {
	// call get category service
	categories, errResponse := service.IMerchantService.GetCategories()
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Data: categories,
	})
}

// get the avilable brands
//
//	@Summary		Get brands
//	@Description	Get all brands
//	@ID				get_brands
//	@Tags			Merchant
//	@Produce		json
//	@Security 	    JWT
//	@Success		200		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		401		{object}	dto.ResponseJson
//	@Failure		403		{object}	dto.ResponseJson
//	@Router			/merchant/brand [get]
func (service *MerchantHandler) GetBrands(ctx *fiber.Ctx) error {
	// call get brand service
	brands, errResponse := service.IMerchantService.GetBrands()
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Data: brands,
	})
}

// get all the products for the merchant
//
//	@Summary		Get products
//	@Description	Get all products of the merchant
//	@ID				get_products_merchant
//	@Tags			Merchant
//	@Produce		json
//	@Security 	    JWT
//	@Param			limit query string false "Enter limit"
//	@Param			offset query string false "Enter offset"
//	@Param			category_name query string false "Enter category_name"
//	@Param			brand_name query string false "Enter brand_name"
//	@Param			price query string false "Enter price"
//	@Param			rating query string false "Enter rating"
//	@Success		200		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		401		{object}	dto.ResponseJson
//	@Failure		403		{object}	dto.ResponseJson
//	@Router			/merchant/product [get]
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

// get a single product by id
//
//	@Summary		Get product
//	@Description	Get product of the merchant
//	@ID				get_product_merchant
//	@Tags			Merchant
//	@Produce		json
//	@Security 	    JWT
//	@Param			product_id path string true "Enter product id"
//	@Success		200		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		404		{object}	dto.ResponseJson
//	@Failure		401		{object}	dto.ResponseJson
//	@Failure		403		{object}	dto.ResponseJson
//	@Router			/merchant/product/{product_id} [get]
func (service *MerchantHandler) GetProduct(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(uuid.UUID)
	id := ctx.Params("id")

	// call get product service
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

// get all the orders for the merchant
//
//	@Summary		Get orders
//	@Description	Get orders for the merchant
//	@ID				get_orders_merchant
//	@Tags			Merchant
//	@Accept			json
//	@Produce		json
//	@Security 	    JWT
//	@Param			limit query string false "Enter limit"
//	@Param			offset query string false "Enter offset"
//	@Param			from_date query string false "Enter from_date"
//	@Param			to_date query string false "Enter to_date"
//	@Success		200		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		401		{object}	dto.ResponseJson
//	@Failure		403		{object}	dto.ResponseJson
//	@Router			/merchant/order [get]
func (service *MerchantHandler) GetOrders(ctx *fiber.Ctx) error {
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

	// call get orders service
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

// get a single order by id
//
//	@Summary		Get order
//	@Description	Get order for the merchant
//	@ID				get_order_merchant
//	@Tags			Merchant
//	@Produce		json
//	@Security 	    JWT
//	@Param			order_id path string true "Enter order id"
//	@Success		200		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		404		{object}	dto.ResponseJson
//	@Failure		401		{object}	dto.ResponseJson
//	@Failure		403		{object}	dto.ResponseJson
//	@Router			/merchant/order/{order_id} [get]
func (service *MerchantHandler) GetOrder(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(uuid.UUID)
	id := ctx.Params("id")

	// call get order service
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

// update product details
//
//	@Summary		Update product
//	@Description	Update product of the merchant
//	@ID				update_product_merchant
//	@Tags			Merchant
//	@Accept			json
//	@Produce		json
//	@Security 	    JWT
//	@Param			product_id path string true "Enter product id"
//	@Param			ProductDetails body models.Products true "Enter Product details"
//	@Success		200		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		404		{object}	dto.ResponseJson
//	@Failure		401		{object}	dto.ResponseJson
//	@Failure		403		{object}	dto.ResponseJson
//	@Router			/merchant/product/{product_id} [patch]
func (service *MerchantHandler) UpdateProduct(ctx *fiber.Ctx) error {
	var product models.Products
	userId := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&product); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	// call update product service
	errResponse := service.IMerchantService.UpdateProduct(userId, &product)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "product updated successfully",
	})
}

// update order_item status
//
//	@Summary		Update order_item
//	@Description	Update order_item status for the merchant
//	@ID				update_order_item_merchant
//	@Tags			Merchant
//	@Accept			json
//	@Produce		json
//	@Security 	    JWT
//	@Param			order_item_id path string true "Enter order_item id"
//	@Param			OrderItemStatus body models.OrderedItems true "Enter OrderItem status"
//	@Success		200		{object}	dto.ResponseJson
//	@Failure		304		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		404		{object}	dto.ResponseJson
//	@Failure		401		{object}	dto.ResponseJson
//	@Failure		403		{object}	dto.ResponseJson
//	@Router			/merchant/order/{order_item_id} [patch]
func (service *MerchantHandler) UpdateOrderStatus(ctx *fiber.Ctx) error {
	var orderItem models.OrderedItems

	id := ctx.Params("id")
	userId := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&orderItem); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	// call update order status service
	errResponse := service.IMerchantService.UpdateOrderStatus(userId, id, orderItem.Status)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "order_item status updated successfully",
	})
}

// update merchant details
//
//	@Summary		Update merchant
//	@Description	Update details of merchant
//	@ID				update_details_merchant
//	@Tags			Merchant
//	@Accept			json
//	@Produce		json
//	@Security 	    JWT
//	@Param			MerchantDetails body models.Users true "Enter merchant details"
//	@Success		200		{object}	dto.ResponseJson
//	@Success		304		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		404		{object}	dto.ResponseJson
//	@Failure		401		{object}	dto.ResponseJson
//	@Failure		403		{object}	dto.ResponseJson
//	@Router			/merchant [patch]
func (service *MerchantHandler) UpdateMerchant(ctx *fiber.Ctx) error {
	var user models.Users
	userId := ctx.Locals("user_id").(uuid.UUID)

	if err := ctx.BodyParser(&user); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	// call update merchant service
	errResponse := service.IMerchantService.UpdateMerchant(userId, &user)
	if errResponse != nil {
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{
			Error: errResponse.Error,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "merchant details updated successfully",
	})
}

// delete product of the merchant
//
//	@Summary		Delete product
//	@Description	Delete product of the merchant
//	@ID				delete_product_merchant
//	@Tags			Merchant
//	@Produce		json
//	@Security 	    JWT
//	@Param			product_id path string true "Enter product id"
//	@Success		200		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		304		{object}	dto.ResponseJson
//	@Failure		404		{object}	dto.ResponseJson
//	@Failure		401		{object}	dto.ResponseJson
//	@Failure		403		{object}	dto.ResponseJson
//	@Router			/merchant/product/{product_id} [delete]
func (service *MerchantHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	// call delete product service
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

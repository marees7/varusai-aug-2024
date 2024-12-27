package repositories

import (
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type IAdminRepository interface {
	AddCategoreyRepository(*models.Categories) *dto.ErrorResponse
	AddBrandRepository(*models.Brands) *dto.ErrorResponse
}

type adminRepository struct {
	*gorm.DB
}

func CommenceAdminRepository(db *gorm.DB) IAdminRepository {
	return &adminRepository{db}
}

func (db *adminRepository) AddCategoreyRepository(category *models.Categories) *dto.ErrorResponse {
	record := db.Where("category_name = ?", category.CategoryName).First(category)
	if record.RowsAffected > 0 {
		return &dto.ErrorResponse{Status: fiber.StatusConflict,
			Error: "category already exists"}
	}

	record = db.Create(category)
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	}

	return nil
}

func (db *adminRepository) AddBrandRepository(brand *models.Brands) *dto.ErrorResponse {
	record := db.Where("category_name = ?", brand.BrandName).First(brand)
	if record.RowsAffected > 0 {
		return &dto.ErrorResponse{Status: fiber.StatusConflict,
			Error: "brand already exists"}
	}

	record = db.Create(brand)
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusInternalServerError,
			Error: record.Error.Error()}
	}

	return nil
}

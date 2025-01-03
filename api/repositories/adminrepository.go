package repositories

import (
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type IAdminRepository interface {
	CreateCategory(*models.Categories) *dto.ErrorResponse
	CreateBrand(*models.Brands) *dto.ErrorResponse
}

type adminRepository struct {
	*gorm.DB
}

func CommenceAdminRepository(db *gorm.DB) IAdminRepository {
	return &adminRepository{db}
}

// create category
func (db *adminRepository) CreateCategory(category *models.Categories) *dto.ErrorResponse {
	// check is the category already avilable
	record := db.Where("category_name = ?", category.CategoryName).First(category)
	if record.RowsAffected > 0 {
		return &dto.ErrorResponse{Status: fiber.StatusConflict,
			Error: "category already exists"}
	}

	// create category
	record = db.Create(category)
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	}

	return nil
}

// create brand
func (db *adminRepository) CreateBrand(brand *models.Brands) *dto.ErrorResponse {
	// check is the brand already avilable
	record := db.Where("category_name = ?", brand.BrandName).First(brand)
	if record.RowsAffected > 0 {
		return &dto.ErrorResponse{Status: fiber.StatusConflict,
			Error: "brand already exists"}
	}

	// create brand
	record = db.Create(brand)
	if record.Error != nil {
		return &dto.ErrorResponse{Status: fiber.StatusBadRequest,
			Error: record.Error.Error()}
	}

	return nil
}

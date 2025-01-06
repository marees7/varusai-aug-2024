package services

import (
	"shopping-site/api/repositories"
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"

	"github.com/gofiber/fiber/v2"
)

type IAdminService interface {
	CreateCategorey(*models.Categories) *dto.ErrorResponse
	CreateBrand(*models.Brands) *dto.ErrorResponse
}

type adminService struct {
	repositories.IAdminRepository
}

func CommenceAdminService(admin repositories.IAdminRepository) IAdminService {
	return &adminService{admin}
}

// create category
func (repo *adminService) CreateCategorey(category *models.Categories) *dto.ErrorResponse {
	if category.CategoryName == "" {
		return &dto.ErrorResponse{
			Status: fiber.StatusBadRequest,
			Error:  "category name should not be empty"}
	}
	return repo.IAdminRepository.CreateCategory(category)
}

// create brand
func (repo *adminService) CreateBrand(brand *models.Brands) *dto.ErrorResponse {
	if brand.BrandName == "" {
		return &dto.ErrorResponse{
			Status: fiber.StatusBadRequest,
			Error:  "brand name should not be empty"}
	}
	return repo.IAdminRepository.CreateBrand(brand)
}

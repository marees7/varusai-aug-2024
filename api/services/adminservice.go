package services

import (
	"shopping-site/api/repositories"
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"
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
	return repo.IAdminRepository.CreateCategory(category)
}

// create brand
func (repo *adminService) CreateBrand(brand *models.Brands) *dto.ErrorResponse {
	return repo.IAdminRepository.CreateBrand(brand)
}

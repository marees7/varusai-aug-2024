package validation

import (
	"fmt"
	"regexp"
	"shopping-site/pkg/models"
	"shopping-site/utils/constants"
	"shopping-site/utils/dto"
)

func ValidateUser(user models.Users) error {
	if (user.FirstName == "") || (user.LastName == "") {
		return fmt.Errorf("both first and last name is manditory")
	}

	if len(user.FirstName) <= 3 || len(user.FirstName) >= 10 {
		return fmt.Errorf("first name length is below or above the limit")
	}

	if len(user.FirstName) <= 3 || len(user.FirstName) >= 10 {
		return fmt.Errorf("last name length is below or above the limit")
	}

	if len(user.Password) <= 7 || len(user.Password) >= 15 {
		return fmt.Errorf("password length is below or above the limit")
	}

	if len(user.Phone) != 10 {
		return fmt.Errorf("invalid phone number")
	}

	regxEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !regxEmail.MatchString(user.Email) {
		return fmt.Errorf("invalid email format")
	}

	lowerCase := regexp.MustCompile(`[a-z]`)
	upperCase := regexp.MustCompile(`[A-Z]`)
	digit := regexp.MustCompile(`\d`)
	specialChar := regexp.MustCompile(`[!@#$%^&*()]`)

	if !lowerCase.MatchString(user.Password) {
		return fmt.Errorf("weak password.include any lower case character")
	}
	if !upperCase.MatchString(user.Password) {
		return fmt.Errorf("weak password.include any upper case character")
	}
	if !digit.MatchString(user.Password) {
		return fmt.Errorf("weak password.include any digit")
	}
	if !specialChar.MatchString(user.Password) {
		return fmt.Errorf("weak password.include any special character")
	}

	if user.Role != constants.AdminRole && user.Role != constants.MerchantRole && user.Role != constants.UserRole {
		return fmt.Errorf("unauthorized role name")
	}

	return nil
}

func ValidateLogin(loginRequest dto.LoginRequest) error {
	if loginRequest.Email == "" {
		return fmt.Errorf("email field should not be empty")
	}

	if loginRequest.Password == "" {
		return fmt.Errorf("password field should not be empty")
	}

	return nil
}

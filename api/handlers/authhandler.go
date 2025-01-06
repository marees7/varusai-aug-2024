package handlers

import (
	"fmt"
	"os"
	"shopping-site/api/services"
	"shopping-site/api/validation"
	"shopping-site/pkg/loggers"
	"shopping-site/pkg/models"
	"shopping-site/utils/dto"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	services.IAuthService
}

// signup a new (User) or (Merchant) or(Admin)
//
//	@Summary		Signup a new member
//	@Description	Create a new member based on the role
//	@ID				signup
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			Signup	body		models.Users	true	"Enter your details"
//	@Success		201		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		409		{object}	dto.ResponseJson
//	@Router			/signup [post]
func (handler *AuthHandler) Signup(ctx *fiber.Ctx) error {
	var user models.Users

	if err := ctx.BodyParser(&user); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{Error: err.Error()})
	}

	// validate the user details
	if err := validation.ValidateUser(user); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	// call the signup service
	if err := handler.IAuthService.SignUp(user); err != nil {
		loggers.ErrorLog.Println(err.Error)
		return ctx.Status(err.Status).JSON(dto.ResponseJson{Error: err.Error})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{
		Message: fmt.Sprintf("%s created successfully", user.Role)})
}

// Login (User) or (Merchant) or(Admin)
//
//	@Summary		Login members
//	@Description	athunticate user and generate token
//	@ID				login
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			Login	body		dto.LoginRequest	true	"Enter login details"
//	@Success		200		{object}	dto.ResponseJson
//	@Failure		400		{object}	dto.ResponseJson
//	@Failure		404		{object}	dto.ResponseJson
//	@Router			/login [post]
func (handler *AuthHandler) Login(ctx *fiber.Ctx) error {
	var loginRequest dto.LoginRequest

	if err := ctx.BodyParser(&loginRequest); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{Error: err.Error()})
	}

	// validate the credentials
	if err := validation.ValidateLogin(loginRequest); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	// call the login service
	user, errResponse := handler.IAuthService.Login(loginRequest)
	if errResponse != nil {
		loggers.ErrorLog.Println(errResponse.Error)
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{Error: errResponse.Error})
	}
	// set payload on jwt token
	claims := &dto.JWTClaims{
		UserID: user.UserID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the jwt token with signing method
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{Error: err.Error()})
	}

	// initiate cookies struct with value
	cookie := fiber.Cookie{
		Name:     "jwt_" + user.Role,
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	// set the cookies
	ctx.Cookie(&cookie)
	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseJson{
		Message: "logged in successfully"})
}

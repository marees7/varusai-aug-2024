package handlers

import (
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

func (handler *AuthHandler) SignupHandler(ctx *fiber.Ctx) error {
	var user models.Users

	if err := ctx.BodyParser(&user); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{Error: err.Error()})
	}

	if err := validation.ValidateUser(user); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	if err := handler.IAuthService.SignUpService(user); err != nil {
		loggers.ErrorLog.Println(err.Error)
		return ctx.Status(err.Status).JSON(dto.ResponseJson{Error: err.Error})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{Message: "User Created"})
}

func (handler *AuthHandler) LoginHandler(ctx *fiber.Ctx) error {
	var loginRequest dto.LoginRequest

	if err := ctx.BodyParser(&loginRequest); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{Error: err.Error()})
	}

	if err := validation.ValidateLogin(loginRequest); err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{
			Error: err.Error(),
		})
	}

	user, errResponse := handler.IAuthService.LoginService(loginRequest)
	if errResponse != nil {
		loggers.ErrorLog.Println(errResponse.Error)
		return ctx.Status(errResponse.Status).JSON(dto.ResponseJson{Error: errResponse.Error})
	}

	claims := &dto.JWTClaims{
		UserID: user.UserId,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseJson{Error: err.Error()})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	ctx.Cookie(&cookie)
	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseJson{Message: "Logged in successfully"})

}

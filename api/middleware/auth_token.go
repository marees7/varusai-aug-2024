package middleware

import (
	"fmt"
	"os"
	"shopping-site/pkg/loggers"
	"shopping-site/utils/dto"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJwt(ctx *fiber.Ctx) error {
	tokenString := ctx.Cookies("jwt")
	if tokenString == "" {
		loggers.WarnLog.Println("Login required to proceed")
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.ResponseJson{
			Message: "Login required to proceed",
			Error:   "unauthorized request",
		})
	}

	token, err := jwt.ParseWithClaims(tokenString, &dto.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.ResponseJson{
			Message: "invalid token",
			Error:   err.Error(),
		})
	}

	claims, ok := token.Claims.(*dto.JWTClaims)
	if !ok {
		loggers.WarnLog.Println(err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.ResponseJson{
			Error: "invalid token",
		})
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		fmt.Println("The claims exp ", claims.ExpiresAt.Unix())
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.ResponseJson{
			Message: "session expired,please login again to proceed",
		})
	}

	ctx.Locals("user_id", claims.UserID)
	ctx.Locals("role", claims.Role)
	ctx.Locals("email", claims.Email)

	return ctx.Next()
}

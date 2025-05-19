package utils

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("SECRET_KAMU_YANG_AMAN") // ganti dengan secret yang aman dan simpan di env

func CreateJWT(userID string, userPublicKey string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":         userID,
		"user_public_key": userPublicKey,
		"exp":             time.Now().Add(time.Hour * 24).Unix(), // token berlaku 1 hari
		"iat":             time.Now().Unix(),                     // waktu pembuatan
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing or invalid token",
			})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}

		// simpan user_id ke context (jika diperlukan)
		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["user_id"])

		return c.Next()
	}
}

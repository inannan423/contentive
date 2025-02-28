package middlewares

import (
	"contentive/config"
	"contentive/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware is a middleware that validates the JWT token and sets the user in the context
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header from the request
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No Authorization header provided",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// If the token string is the same as the original header, it means the prefix was not found
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}

		// Validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// turn token to jwt.MapClaims
		claims := token.Claims.(jwt.MapClaims)
		// get user_id from claims
		userID := claims["user_id"].(string)

		var user models.User
		if err := config.DB.Preload("Role.Permissions").First(&user, "id = ?", userID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Use not found",
			})
		}

		// Set the user in the context
		c.Locals("user", &user)
		return c.Next() // continue to the next handler
	}
}

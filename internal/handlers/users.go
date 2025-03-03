package handlers

import (
	"contentive/config"
	"contentive/internal/models"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Login: Login handler
// Input: Email, Password
// Output: Token, User Details
func Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Cannot find user
	var user models.User
	if err := config.DB.Preload("Role.Permissions").Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if !user.CheckPassword(input.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	// Set expiration time for the token
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	user.LastLogin = time.Now()
	config.DB.Save(&user)

	return c.JSON(fiber.Map{
		"token": t,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role.Type,
		},
	})
}

// CreateUser: Create user handler
// Input: Username, Email, Password, RoleID
// Output: User Details
func CreateUser(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		RoleID   string `json:"role_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	roleID, err := uuid.Parse(input.RoleID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role ID"})
	}

	// Check if the role is super_admin
	var role models.Role
	if err := config.DB.First(&role, roleID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role"})
	}

	if role.Type == "super_admin" {
		// Check if a super_admin user already exists
		var count int64
		if err := config.DB.Model(&models.User{}).Joins("JOIN roles ON users.role_id = roles.id").Where("roles.type = ?", "super_admin").Count(&count).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check existing super admin"})
		}

		if count > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Super admin user already exists"})
		}
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		RoleID:   roleID,
	}

	// Hash the password
	if err := user.HashPassword(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Save the user to the database
	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	if err := config.DB.Preload("Role.Permissions").First(&user, user.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load user data"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role_id":    user.RoleID,
		"role":       user.Role,
		"active":     user.Active,
		"last_login": user.LastLogin,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	if err := config.DB.Preload("Role.Permissions").Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	var safeUsers []fiber.Map
	for _, user := range users {
		safeUsers = append(safeUsers, fiber.Map{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"role_id":    user.RoleID,
			"role":       user.Role,
			"active":     user.Active,
			"last_login": user.LastLogin,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}

	return c.JSON(safeUsers)
}

func UpdateUser(c *fiber.Ctx) error {
	user := c.Locals("targetUser").(*models.User)

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		RoleID   string `json:"role_id"`
		Active   *bool  `json:"active"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	log.Printf("Update input: %+v", input)

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Password != "" {
		user.Password = input.Password
		if err := user.HashPassword(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
		}
	}
	if input.RoleID != "" {
		roleID, err := uuid.Parse(input.RoleID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role ID"})
		}

		var role models.Role
		if err := config.DB.Preload("Permissions").First(&role, "id = ?", roleID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Role not found"})
		}

		// Check if trying to update to super_admin role
		if role.Type == models.SuperAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Cannot update user to super admin role"})
		}

		user.RoleID = roleID
		user.Role = role
	}
	if input.Active != nil {
		user.Active = *input.Active
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	if err := config.DB.Preload("Role.Permissions").First(&user, user.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load updated user data"})
	}

	return c.JSON(fiber.Map{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role_id":    user.RoleID,
		"role":       user.Role,
		"active":     user.Active,
		"last_login": user.LastLogin,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user models.User
	if err := config.DB.Preload("Role").First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if user.Role.Type == models.SuperAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot delete super admin user",
		})
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func ValidateToken(c *fiber.Ctx) error {
	// If can reach here, it means the token is valid
	return c.JSON(fiber.Map{
		"valid": true,
	})
}

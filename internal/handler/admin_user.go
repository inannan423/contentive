package handler

import (
	"contentive/internal/database"
	"contentive/internal/logger"
	"contentive/internal/models"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetAllAdminUsers returns all admin users
func GetAllAdminUsers(c *fiber.Ctx) error {
	var users []models.AdminUser

	if err := database.DB.Find(&users).Error; err != nil {
		logger.Error("Failed to fetch users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	currentUser := c.Locals("user").(models.AdminUser)
	logger.AdminAction(
		currentUser.ID,
		currentUser.Name,
		"GET_ALL_USERS",
		"Retrieved all users list",
	)

	return c.Status(fiber.StatusOK).JSON(users)
}

// GetAdminUserById returns a admin user by id
func GetAdminUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.AdminUser
	if err := database.DB.First(&user, id).Error; err != nil {
		logger.Error("Failed to fetch user: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	currentUser := c.Locals("user").(models.AdminUser)
	logger.AdminAction(
		currentUser.ID,
		currentUser.Name,
		"GET_USER_BY_ID",
		"Retrieved user with id: "+id,
	)

	return c.Status(fiber.StatusOK).JSON(user)
}

// CreateAdminUser creates a new admin user
func CreateAdminUser(c *fiber.Ctx) error {
	var input struct {
		Name     string                 `json:"name"`
		Email    string                 `json:"email"`
		Password string                 `json:"password"`
		Role     models.AdminUserRole   `json:"role"`
		Status   models.AdminUserStatus `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		logger.Error("Failed to parse input: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Check if the required fields are present
	if input.Name == "" || input.Email == "" || input.Password == "" || input.Role == "" || input.Status == "" {
		logger.Error("Missing required fields")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: name, email, password, role, status",
		})
	}

	// Check if the password is at least 8 characters long and contains at least one uppercase letter, one lowercase letter, and one number
	if len(input.Password) < 8 || !containsUppercase(input.Password) || !containsLowercase(input.Password) || !containsNumber(input.Password) {
		logger.Error("Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one number")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one number",
		})
	}

	// Check if the email is valid
	if !isValidEmail(input.Email) {
		logger.Error("Invalid email")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email",
		})
	}

	// Check if the role is valid
	if input.Role != models.AdminUserRoleViewer && input.Role != models.AdminUserRoleEditor && input.Role != models.AdminUserRoleAdmin {
		logger.Error("Invalid role")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role, must be one of: viewer, editor, admin",
		})
	}

	// Check if the status is valid
	if input.Status != models.AdminUserStatusActive && input.Status != models.AdminUserStatusInactive {
		logger.Error("Invalid status")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status, must be one of: active, inactive",
		})
	}

	// Check if the name or email already exists
	var existingUser models.AdminUser
	if err := database.DB.Where("name = ? OR email = ?", input.Name, input.Email).First(&existingUser).Error; err == nil {
		logger.Error("Name or email already exists")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name or email already exists",
		})
	}
	// Create the new user
	user := models.AdminUser{
		Name:        input.Name,
		Email:       input.Email,
		Password:    input.Password,
		Role:        input.Role,
		Status:      input.Status,
		LastLoginAt: time.Now(),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		logger.Error("Failed to create user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	currentUser := c.Locals("user").(models.AdminUser)
	logger.AdminAction(
		currentUser.ID,
		currentUser.Name,
		"CREATE_USER",
		"Created new user: "+user.Name,
	)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"status":     user.Status,
		"created_at": user.CreatedAt,
	})
}

// UpdateAdminUser updates an admin user
func UpdateAdminUser(c *fiber.Ctx) error {
	id := c.Params("id")

	currentUser := c.Locals("user").(models.AdminUser)

	var user models.AdminUser
	if err := database.DB.First(&user, id).Error; err != nil {
		logger.Error("Failed to fetch user: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if user.Role == models.AdminUserRoleSuperAdmin && currentUser.Role != models.AdminUserRoleSuperAdmin {
		logger.Error("Attempted to modify super admin user by non-super admin user")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot modify super admin user",
		})
	}

	var input struct {
		Name     *string                 `json:"name,omitempty"`
		Email    *string                 `json:"email,omitempty"`
		Password *string                 `json:"password,omitempty"`
		Role     *models.AdminUserRole   `json:"role,omitempty"`
		Status   *models.AdminUserStatus `json:"status,omitempty"`
	}

	if err := c.BodyParser(&input); err != nil {
		logger.Error("Failed to parse input: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	if input.Name != nil {
		if *input.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Name cannot be empty",
			})
		}
		var existingUser models.AdminUser
		if err := database.DB.Where("name = ? AND id != ?", *input.Name, id).First(&existingUser).Error; err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Name already exists",
			})
		}
		user.Name = *input.Name
	}

	if input.Email != nil {
		if *input.Email == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email cannot be empty",
			})
		}
		if !isValidEmail(*input.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid email format",
			})
		}
		var existingUser models.AdminUser
		if err := database.DB.Where("email = ? AND id != ?", *input.Email, id).First(&existingUser).Error; err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		user.Email = *input.Email
	}

	if input.Password != nil {
		if *input.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Password cannot be empty",
			})
		}
		if len(*input.Password) < 8 || !containsUppercase(*input.Password) || !containsLowercase(*input.Password) || !containsNumber(*input.Password) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one number",
			})
		}
		if err := user.SetPassword(*input.Password); err != nil {
			logger.Error("Failed to hash password: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update password",
			})
		}
	}

	if input.Role != nil {
		if *input.Role == models.AdminUserRoleSuperAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cannot set role to super_admin",
			})
		}
		if *input.Role != models.AdminUserRoleViewer && *input.Role != models.AdminUserRoleEditor && *input.Role != models.AdminUserRoleAdmin {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid role",
			})
		}
		user.Role = *input.Role
	}

	if input.Status != nil {
		if *input.Status != models.AdminUserStatusActive && *input.Status != models.AdminUserStatusInactive {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid status",
			})
		}
		user.Status = *input.Status
	}

	if err := database.DB.Save(&user).Error; err != nil {
		logger.Error("Failed to update user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	logger.AdminAction(
		currentUser.ID,
		currentUser.Name,
		"UPDATE_USER",
		"Updated user: "+user.Name,
	)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"status":     user.Status,
		"updated_at": user.UpdatedAt,
	})
}

// DeleteAdminUser deletes an admin user
func DeleteAdminUser(c *fiber.Ctx) error {
	id := c.Params("id")

	currentUser := c.Locals("user").(models.AdminUser)

	var user models.AdminUser
	if err := database.DB.First(&user, id).Error; err != nil {
		logger.Error("Failed to fetch user: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if the user is trying to delete themselves
	if user.ID == currentUser.ID {
		logger.Error("Attempted to delete self")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot delete self",
		})
	}

	// Check if the user is trying to delete a super admin user
	if user.Role == models.AdminUserRoleSuperAdmin {
		logger.Error("Attempted to delete super admin user")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot delete super admin user",
		})
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		logger.Error("Failed to delete user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	logger.AdminAction(
		currentUser.ID,
		currentUser.Name,
		"DELETE_USER",
		"Deleted user: "+user.Name,
	)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func containsUppercase(s string) bool {
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			return true
		}
	}
	return false
}

func containsLowercase(s string) bool {
	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			return true
		}
	}
	return false
}

func containsNumber(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

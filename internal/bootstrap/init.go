package bootstrap

import (
	"contentive/config"
	"contentive/internal/models"
	"log"
)

func InitRolesAndPermissions() {
	permissions := []models.Permission{
		// Permission for creating content types
		{
			Name:        "Create Content Type",
			Type:        models.CreateContentType,
			Description: "Roles with this permission can create new content types",
		},
		{
			Name:        "Read Content Type",
			Type:        models.ReadContentType,
			Description: "Roles with this permission can read content types",
		},
		{
			Name:        "Update Content Type",
			Type:        models.UpdateContentType,
			Description: "Roles with this permission can update content types",
		},

		// Permission for creating content
		{
			Name:        "Create Content",
			Type:        models.CreateContent,
			Description: "Roles with this permission can create new content",
		},
		{
			Name:        "Read Content",
			Type:        models.ReadContent,
			Description: "Roles with this permission can read content",
		},
		{
			Name:        "Update Content",
			Type:        models.UpdateContent,
			Description: "Roles with this permission can update content",
		},
		{
			Name:        "Delete Content",
			Type:        models.DeleteContent,
			Description: "Roles with this permission can delete content",
		},

		// Permission for managing users
		{
			Name:        "Manage Users",
			Type:        models.ManageUsers,
			Description: "Roles with this permission can manage users",
		},
		{
			Name:        "Manage Roles",
			Type:        models.ManageRoles,
			Description: "Roles with this permission can manage roles",
		},
		{
			Name:        "View Audit Logs",
			Type:        models.ViewAuditLogs,
			Description: "Roles with this permission can view audit logs",
		},
	}

	for _, p := range permissions {
		if err := config.DB.Where(models.Permission{Type: p.Type}).FirstOrCreate(&p).Error; err != nil {
			log.Printf("Error creating permission %s: %v", p.Name, err)
		}
	}

	// Get all permissions
	var allPermissions []models.Permission
	config.DB.Find(&allPermissions)

	roles := []models.Role{
		{
			Name:        "Super Admin",
			Type:        models.SuperAdmin,
			Description: "Has full access to all features",
			Permissions: allPermissions,
		},
		{
			Name:        "Content Admin",
			Type:        models.ContentAdmin,
			Description: "Can manage all content and content types",
			Permissions: filterPermissions(allPermissions, []models.PermissionType{
				models.CreateContentType, models.ReadContentType, models.UpdateContentType,
				models.CreateContent, models.ReadContent, models.UpdateContent, models.DeleteContent,
				models.ViewAuditLogs,
			}),
		},
		{
			Name:        "Editor",
			Type:        models.Editor,
			Description: "Can create and edit content",
			Permissions: filterPermissions(allPermissions, []models.PermissionType{
				models.ReadContentType,
				models.CreateContent, models.ReadContent, models.UpdateContent,
			}),
		},
		{
			Name:        "Viewer",
			Type:        models.Viewer,
			Description: "Can only view content",
			Permissions: filterPermissions(allPermissions, []models.PermissionType{
				models.ReadContentType,
				models.ReadContent,
			}),
		},
	}

	for _, r := range roles {
		var existingRole models.Role
		if err := config.DB.Where(models.Role{Type: r.Type}).FirstOrCreate(&existingRole).Error; err != nil {
			log.Printf("Error creating role %s: %v", r.Name, err)
			continue
		}

		if err := config.DB.Model(&existingRole).Association("Permissions").Replace(r.Permissions); err != nil {
			log.Printf("Error setting permissions for role %s: %v", r.Name, err)
		}
	}
}

// filterPermissions filters the given permissions based on the given types
func filterPermissions(allPermissions []models.Permission, types []models.PermissionType) []models.Permission {
	var filtered []models.Permission
	for _, p := range allPermissions {
		for _, t := range types {
			if p.Type == t {
				filtered = append(filtered, p)
				break
			}
		}
	}
	return filtered
}

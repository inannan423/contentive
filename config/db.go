package config

import (
	"contentive/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=" + AppConfig.DBHost + " user=" + AppConfig.DBUser + " dbname=" + AppConfig.DBName + " password=" + AppConfig.DBPassword + " port=" + AppConfig.DBPort + " sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}
	log.Println("Connected to the database!")

	// check if uuid-ossp extension exists
	var hasExtension bool
	err = DB.Raw(`SELECT EXISTS (
        SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp'
    )`).Scan(&hasExtension).Error
	if err != nil {
		log.Fatal("failed to check uuid-ossp extension:", err)
	}
	log.Printf("UUID extension exists: %v\n", hasExtension)

	// if uuid-ossp extension does not exist, create it
	if !hasExtension {
		var isSuperuser bool
		if err := DB.Raw(`SELECT rolsuper FROM pg_roles WHERE rolname = current_user`).Scan(&isSuperuser).Error; err != nil {
			log.Fatal("failed to check user privileges:", err)
		}
		log.Printf("Current user is superuser: %v\n", isSuperuser)

		if err := DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
			log.Printf("Error creating uuid-ossp extension: %v\n", err)
			if !isSuperuser {
				log.Fatal("Current user does not have superuser privileges. Please run:\nALTER USER contentive WITH SUPERUSER;")
			}
			log.Fatal("Failed to create uuid-ossp extension")
		}
		log.Println("UUID extension created successfully!")
	}

	// Force drop tables
	// if err := DB.Migrator().DropTable(&models.ContentEntry{}, &models.Field{}, &models.ContentType{}); err != nil {
	// 	log.Fatal("failed to drop tables:", err)
	// }
	// log.Println("Existing tables dropped!")

	// Migrate the schema
	if err := DB.AutoMigrate(
		&models.ContentType{},
		&models.Field{},
		&models.ContentEntry{},
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.AuditLog{},
		&models.APIRole{},
		&models.APIPermission{},
	); err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	log.Println("Database migration completed!")
}

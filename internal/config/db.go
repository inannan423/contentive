package config

import (
	"contentive/internal/logger"
	"contentive/internal/models"
	"fmt"
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

	if err := DB.Exec(`DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'relation_type_enum') THEN
			CREATE TYPE relation_type_enum AS ENUM ('one_to_one', 'one_to_many', 'many_to_one', 'many_to_many');
		END IF;
	END $$;`).Error; err != nil {
		logger.GeneralAction(fmt.Sprintf("Error creating relation_type_enum: %v", err))
		log.Fatal("failed to create relation_type_enum:", err)
	}
	logger.GeneralAction("Enum types created successfully!")

	// check if uuid-ossp extension exists
	var hasExtension bool
	err = DB.Raw(`SELECT EXISTS (
        SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp'
    )`).Scan(&hasExtension).Error
	if err != nil {
		logger.GeneralAction(fmt.Sprintf("Error checking uuid-ossp extension: %v", err))
		log.Fatal("failed to check uuid-ossp extension:", err)
	}
	logger.GeneralAction(fmt.Sprintf("Checking if UUID extension exists: %v", hasExtension))

	// if uuid-ossp extension does not exist, create it
	if !hasExtension {
		var isSuperuser bool
		if err := DB.Raw(`SELECT rolsuper FROM pg_roles WHERE rolname = current_user`).Scan(&isSuperuser).Error; err != nil {
			logger.GeneralAction(fmt.Sprintf("Error checking user privileges: %v", err))
			log.Fatal("failed to check user privileges:", err)
		}

		if err := DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
			logger.GeneralAction(fmt.Sprintf("Error creating uuid-ossp extension: %v", err))
			if !isSuperuser {
				logger.GeneralAction("Current user does not have superuser privileges. Please run:\nALTER USER contentive WITH SUPERUSER;")
				log.Fatal("Current user does not have superuser privileges. Please run:\nALTER USER contentive WITH SUPERUSER;")
			}
			logger.GeneralAction(fmt.Sprintf("Error creating uuid-ossp extension: %v", err))
			log.Fatal("Failed to create uuid-ossp extension")
		}
		logger.GeneralAction("UUID extension created successfully")
	}

	// Force drop tables
	// if err := DB.Migrator().DropTable(&models.ContentEntry{}, &models.Field{}, &models.ContentType{}); err != nil {
	// 	log.Fatal("failed to drop tables:", err)
	// }
	// log.Println("Existing tables dropped!")

	// Migrate the schema
	if err := DB.AutoMigrate(
		&models.AdminUser{},
		&models.APIUser{},
	); err != nil {
		logger.GeneralAction(fmt.Sprintf("Error migrating database: %v", err))
		log.Fatal("failed to migrate database:", err)
	}
	logger.GeneralAction("Database migration completed")
}

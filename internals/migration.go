package internals

import (
	"fmt"
	"shopping-site/pkg/loggers"
	"shopping-site/pkg/models"

	"gorm.io/gorm"
)

func SchemaMigration(db *gorm.DB) {
	err := db.AutoMigrate(&models.Users{}, &models.Addresses{}, &models.Categories{}, &models.Brands{}, &models.Products{}, &models.Orders{}, &models.OrderedItems{})
	if err != nil {
		loggers.FatalLog.Fatal("Error while migrating tables")
	}

	loggers.InfoLog.Print("Migration Completed")
	fmt.Println("Migration Completed")
}

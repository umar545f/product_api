package migrate

import (
	"fmt"
	"go-microservices/product-api/data"

	"gorm.io/gorm"
)

func MigrateAllTables(r *gorm.DB) error {
	err := data.MigrateProduct(r)
	if err != nil {
		fmt.Printf("Could not migrate all table %s", err.Error())
		return err
	}
	return nil
}

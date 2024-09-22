package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

var ErrorProduct = fmt.Errorf("Product Not Found")

// swagger: model Product
type Product struct {
	// The ID for this product
	// required: true
	// min: 1
	ID int `gorm:"primaryKey;autoIncrement" json:"id"`

	// The name for this product
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// The description for this product
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// The price for the product
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"gte=1"`

	// The SKU for the product
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"required,sku"`
}

func MigrateProduct(r *gorm.DB) error {
	err := r.AutoMigrate(&Product{})
	if err != nil {
		fmt.Printf("Could not migrate product table %s", err.Error())
		return err
	}
	return nil

}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)

	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku is of format abc-absd-dfsdf
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

func (p *Product) FromJson(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(p)
	return err
}

func GetProducts(db *gorm.DB) (*[]Product, error) {
	products := &[]Product{}
	err := db.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Adding a new product
func AddProduct(p *Product, db *gorm.DB) error {
	err := db.Create(p).Error
	if err != nil {
		return err
	}
	return nil

}

// Updating the product based on ID
func UpdateProducts(id int, newProd *Product, db *gorm.DB) error {
	product := Product{}
	err := db.First(&product, id).Error
	if err != nil {
		return err
	}

	product.Name = newProd.Name
	product.Description = newProd.Description
	product.Price = newProd.Price
	product.SKU = newProd.SKU

	err = db.Save(&product).Error
	if err != nil {
		return err
	}

	return nil

}

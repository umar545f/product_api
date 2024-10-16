// Package classification of Product API
//
// # Documentation for Product API
//
// Schemes: http
// BasePath: /products
// Version: 1.0.0
//
// Consumes:
// - application/json
// swagger:meta
package handler

import (
	"context"
	"log"
	"net/http"

	"go-microservices/data"

	"gorm.io/gorm"
)

// swagger:response productsResponse
type productsResponseWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// swagger:parameters updateProducts
type productIdParameterWrapper struct {
	//id of the product to be updated
	//in: path
	//required: true
	ID int `json:"id"`
}

// swagger:response updatedSuccessfully
type productsUpdatedWrapper struct {
}

type Products struct {
	l  *log.Logger
	db *gorm.DB
}

// returns the handler for Products which can be used to call different functions
func NewProducts(l *log.Logger, db *gorm.DB) *Products {
	return &Products{l, db}
}

type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	p.l.Println("MiddleWare to validate request body")

	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			prod := &data.Product{}
			//serialize request body json into prod struct
			err := prod.FromJson(r)

			if err != nil {
				p.l.Println("[Error] deserializing product", err.Error())
				http.Error(rw, "Error reading product", http.StatusInternalServerError)
				return
			}

			// validate json data as per the validation we put on Product struct
			//in data/products
			err = prod.Validate()
			if err != nil {
				p.l.Println("[Error] validating json", err.Error())
				http.Error(rw, "Error validating json", http.StatusInternalServerError)
				return
			}
			//add the product to the context
			ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
			r = r.WithContext(ctx)

			//Call the next handler which will be GET,POST OR PUT

			next.ServeHTTP(rw, r)

		})
}

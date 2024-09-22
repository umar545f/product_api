package handler

import (
	"go-microservices/product-api/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// swagger:route PUT /products/{id} products updateProducts
// Return a list of products
// responses:
//
//	204: updatedSuccessfully
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT reuqest")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Invalid product ID", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProducts(id, prod, p.db)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(rw, "Product not found", http.StatusNotFound)
		} else {
			http.Error(rw, "Could not update the product", http.StatusInternalServerError)
		}
		return
	}

	rw.WriteHeader(http.StatusNoContent) // 204 No Content

}

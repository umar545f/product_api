package handler

import (
	"net/http"

	"go-microservices/data"
)

func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST reuqest")
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	err := data.AddProduct(prod, p.db)
	if err != nil {
		http.Error(rw, "Could not add product", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

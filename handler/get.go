package handler

import (
	"encoding/json"
	"go-microservices/product-api/data"
	"net/http"
)

// swagger:route GET /products products listProducts
// Return a list of products
// responses:
//	200: productsResponse

func (p *Products) GetProducts(rw http.ResponseWriter, _ *http.Request) {
	pL, err := data.GetProducts(p.db)
	if err != nil {
		http.Error(rw, "Could not fetch the products", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(rw).Encode(pL)

	if err != nil {
		http.Error(rw, "Oops,could not marshal the json", http.StatusInternalServerError)
		return
	}
}

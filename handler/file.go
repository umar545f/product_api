package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"go-microservices/data"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func (p *Products) UploadFile(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(rw, "Could not parse file", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")

	if err != nil {
		http.Error(rw, "Could not get file from form", http.StatusBadRequest)
		return
	}

	defer file.Close()

	//create a directory for images
	dir := "images"
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		http.Error(rw, "Could notcreate the directory", http.StatusBadRequest)
		return
	}

	//create a new file to save the image
	imagePath := filepath.Join(dir, fmt.Sprintf("%s.jpeg", id)) // Change extension as needed
	out, err := os.Create(imagePath)
	if err != nil {
		http.Error(rw, "Unable to save the image", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	//Copy the uploaded file in the new file

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(rw, "Unable to copy the image", http.StatusInternalServerError)
		return
	}
	err = data.UploadFile(p.db, imagePath, id)
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

func (p *Products) DownloadFile(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	path, err := data.DownloadFile(p.db, id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(rw, "Image path not found", http.StatusNotFound)
		} else {
			http.Error(rw, "Could not fetch the image", http.StatusInternalServerError)
		}
		return
	}

	http.ServeFile(rw, r, path)

}

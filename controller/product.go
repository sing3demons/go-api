package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/api/models"
	"github.com/sing3demons/api/service"
	"github.com/sing3demons/api/utils"
)

type ProductController interface {
	Create(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
	FindOne(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type productController struct {
	service service.ProductService
}

type Product struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
	Image string `json:"image"`
}

func NewProductController(service service.ProductService) productController {
	return productController{service: service}
}

func (s productController) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := s.service.DeleteProduct(id)
	if err != nil {
		utils.JSON(w, http.StatusNotFound)(map[string]interface{}{"error": err})
	}

	utils.JSON(w, http.StatusNoContent)

}

func (s productController) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	product.Name = r.FormValue("name")
	product.Desc = r.FormValue("desc")
	product.Price, _ = strconv.Atoi(r.FormValue("price"))

	resp, err := s.service.Create(&product)
	if err != nil {
		utils.JSON(w, http.StatusUnprocessableEntity)(map[string]interface{}{"error": err})
		return
	}
	s.saveProductImage(w, r, &product)

	utils.JSON(w, http.StatusCreated)(map[string]interface{}{"product": resp})
}

func (s productController) FindAll(w http.ResponseWriter, r *http.Request) {
	products, err := s.service.FindAll()
	if err != nil {
		utils.JSON(w, http.StatusUnprocessableEntity)(map[string]interface{}{"error": err})
		return
	}

	utils.JSON(w, http.StatusOK)(map[string]interface{}{"product": products})
}

func (s productController) FindOne(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	product, err := s.service.FindOne(id)
	if err != nil {
		utils.JSON(w, http.StatusNotFound)(map[string]interface{}{"error": err})
		return
	}

	utils.JSON(w, http.StatusOK)(map[string]interface{}{"product": product})
}

func (s productController) Update(w http.ResponseWriter, r *http.Request) {
	var form models.Product
	form.Name = r.FormValue("name")
	form.Desc = r.FormValue("desc")
	form.Price, _ = strconv.Atoi(r.FormValue("price"))

	id := mux.Vars(r)["id"]

	product, err := s.service.FindOne(id)
	if err != nil {
		utils.JSON(w, http.StatusNotFound)(map[string]interface{}{"error": err})
		return
	}

	form.ID = product.ID

	if form.Image == "" {
		form.Image = product.Image
	}

	copier.Copy(&product, &form)

	fmt.Printf("id: %v product id %v", id, product.ID)
	s.service.Update(product)
	s.saveProductImage(w, r, product)

	utils.JSON(w, http.StatusOK)(map[string]interface{}{"product": product})
}

func (p *productController) checkProduckImage(image string) {
	if image != "" {
		image = strings.Replace(image, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + image)
	}
}

func (p *productController) saveProductImage(w http.ResponseWriter, r *http.Request, product *models.Product) error {
	file, handler, err := r.FormFile("image")
	// truncated for brevity

	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	if err != nil || file == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	defer file.Close()

	p.checkProduckImage(product.Image)

	path := "uploads/products/" + strconv.Itoa(int(product.ID))

	// Create the uploads folder if it doesn't
	// already exist
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	filename := path + "/" + fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
	product.Image = os.Getenv("HOST") + "/" + filename

	// Create a new file in the uploads directory
	// dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename)))
	dst, err := os.Create(filename)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	p.service.SaveFile(product)

	// fmt.Fprintf(w, "Upload successful\n")

	return nil
}

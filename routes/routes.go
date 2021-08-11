package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sing3demons/api/controller"
	"github.com/sing3demons/api/database"
	"github.com/sing3demons/api/repository"
	"github.com/sing3demons/api/service"
)

func Serve(r *mux.Router) {
	db := database.InitDatabase()
	rdb := database.NewRedisCache("redis:6379", 1, 10)
	v1 := "/api/v1"

	productRepository := repository.NewProductRepository(db, rdb)
	productService := service.NewProductService(productRepository)
	productController := controller.NewProductController(productService)

	productsGroup := fmt.Sprintf(v1 + "/products")

	{
		r.HandleFunc(productsGroup, productController.Create).Methods(http.MethodPost)
		r.HandleFunc(productsGroup, productController.FindAll).Methods(http.MethodGet)
		r.HandleFunc(productsGroup+"/{id}", productController.FindOne).Methods(http.MethodGet)
		r.HandleFunc(productsGroup+"/{id}", productController.Update).Methods(http.MethodPut)
		r.HandleFunc(productsGroup+"/{id}", productController.Delete).Methods(http.MethodDelete)
	}

}

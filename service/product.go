package service

import (
	"log"

	"github.com/sing3demons/api/database"
	"github.com/sing3demons/api/models"
	"github.com/sing3demons/api/repository"
)

type ProductService interface {
	Create(product *models.Product) (*models.Product, error)
	FindAll() ([]models.Product, error)
	FindOne(id string) (*models.Product, error)
	SaveFile(product *models.Product) (*models.Product, error)
	Update(product *models.Product) (*models.Product, error)
	DeleteProduct(id string) error
}

type productService struct {
	repository repository.ProductRepository
	rdb        database.RedisCache
}

func NewProductService(repository repository.ProductRepository) ProductService {
	return &productService{
		repository: repository,
	}
}

func (service productService) DeleteProduct(id string) error {
	err := service.repository.Delete(id)

	if err != nil {
		log.Printf("Failed map %v: ", err)
		return err
	}
	return nil
}

func (service productService) Create(product *models.Product) (*models.Product, error) {
	response, err := service.repository.Create(product)
	if err != nil {
		log.Printf("Failed map %v: ", err)
		return nil, err
	}
	return response, nil
}

func (service productService) FindAll() ([]models.Product, error) {
	products, err := service.repository.GetProducts()
	// products, err
	if err != nil {
		log.Printf("Failed map %v: ", err)
		return nil, err
	}

	return products, nil
}

func (service productService) FindOne(id string) (*models.Product, error) {
	product, err := service.repository.GetProduct(id)
	if err != nil {
		log.Printf("not found %v: ", err)
		return nil, err
	}

	return product, nil
}

func (service productService) SaveFile(product *models.Product) (*models.Product, error) {
	response, err := service.repository.SaveFile(product)
	if err != nil {
		log.Printf("Failed map %v: ", err)
		return nil, err
	}
	return response, nil
}

func (service productService) Update(product *models.Product) (*models.Product, error) {
	response, err := service.repository.Update(product)
	if err != nil {
		log.Printf("Failed map %v: ", err)
		return nil, err
	}
	return response, nil
}

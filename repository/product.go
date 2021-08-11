package repository

import (
	"fmt"
	"log"

	"github.com/sing3demons/api/database"
	"github.com/sing3demons/api/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProducts() ([]models.Product, error)
	Create(product *models.Product) (*models.Product, error)
	GetProduct(string) (*models.Product, error)
	SaveFile(product *models.Product) (*models.Product, error)
	Update(product *models.Product) (*models.Product, error)
	Delete(id string) error
}

type productRepository struct {
	db  *gorm.DB
	rdb database.RedisCache
}

func NewProductRepository(db *gorm.DB, rdb database.RedisCache) ProductRepository {
	return &productRepository{db: db, rdb: rdb}
}

func (r productRepository) GetProducts() ([]models.Product, error) {
	rProduct, _ := r.rdb.Get("products")

	var product []models.Product = rProduct

	if product != nil {
		fmt.Println("Get...Redis")
		product, err := r.rdb.Get("products")
		if err != nil {
			log.Printf("get product :%v\n", err)
		}

		return product, nil
	}

	if err := r.db.Find(&product).Error; err != nil {
		return nil, err
	}

	r.rdb.Set("products", product)

	return product, nil
}

func (r productRepository) Create(product *models.Product) (*models.Product, error) {
	if err := r.db.Create(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r productRepository) Delete(id string) error {
	var product models.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return err
	}

	if err := r.db.Unscoped().Delete(&product).Error; err != nil {
		return err
	}

	return nil
}

func (r productRepository) GetProduct(id string) (*models.Product, error) {
	var product models.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r productRepository) Update(product *models.Product) (*models.Product, error) {
	if err := r.db.Model(&models.Product{}).Updates(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r productRepository) SaveFile(product *models.Product) (*models.Product, error) {
	if err := r.db.Save(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

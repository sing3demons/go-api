package repository

import (
	"fmt"

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
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
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

	fmt.Printf("product: %v\n", product)

	if err := r.db.Delete(&product).Error; err != nil {
		return err
	}

	return nil
}

func (r productRepository) GetProducts() ([]models.Product, error) {
	var product []models.Product
	if err := r.db.Find(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r productRepository) GetProduct(id string) (*models.Product, error) {
	var product models.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r productRepository) SaveFile(product *models.Product) (*models.Product, error) {
	if err := r.db.Save(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r productRepository) Update(product *models.Product) (*models.Product, error) {
	if err := r.db.Model(&product).Updates(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

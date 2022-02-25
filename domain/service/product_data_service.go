package service

import (
	"github.com/kiven-man/product/domain/model"
	"github.com/kiven-man/product/domain/repository"
)

type IProductDataService interface {
	AddProduct(*model.Product) (int64, error)
	DeleteProduct(int64) error
	UpdateProduct(*model.Product) error
	FindProductById(int64) (*model.Product, error)
	FindAllProduct() ([]model.Product, error)
}

func NewProductDataService(productRepository repository.IProductRepository) IProductDataService {
	return &ProductDataService{ProductRepository: productRepository}
}

type ProductDataService struct {
	ProductRepository repository.IProductRepository
}

func (p *ProductDataService) AddProduct(product *model.Product) (int64, error) {
	return p.ProductRepository.CreateProduct(product)
}
func (p *ProductDataService) DeleteProduct(id int64) error {
	return p.ProductRepository.DeleteProduct(id)
}
func (p *ProductDataService) UpdateProduct(product *model.Product) error {
	return p.ProductRepository.UpdateProduct(product)
}
func (p *ProductDataService) FindProductById(id int64) (*model.Product, error) {
	return p.ProductRepository.FindProductById(id)
}
func (p *ProductDataService) FindAllProduct() ([]model.Product, error) {
	return p.ProductRepository.FindAllProduct()
}

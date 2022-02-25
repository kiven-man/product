package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kiven-man/product/domain/model"
)

type IProductRepository interface {
	InitTable() error
	CreateProduct(*model.Product) (int64, error)
	DeleteProduct(int64) error
	UpdateProduct(*model.Product) error
	FindProductById(int64) (*model.Product, error)
	FindAllProduct() ([]model.Product, error)
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{
		mysqlDB: db,
	}
}

type ProductRepository struct {
	mysqlDB *gorm.DB
}

func (p *ProductRepository) InitTable() error {
	return p.mysqlDB.CreateTable(&model.Product{}, &model.ProductSeo{}, &model.ProductImage{}, &model.ProductSize{}).Error
}

func (p *ProductRepository) CreateProduct(product *model.Product) (int64, error) {
	return product.ID, p.mysqlDB.Create(product).Error
}

func (p *ProductRepository) DeleteProduct(id int64) error {
	// 开启事务
	tx := p.mysqlDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	// 删除
	if err := tx.Unscoped().Where("id = ?", id).Delete(&model.Product{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Unscoped().Where("image_product_id = ? ", id).Delete(&model.ProductImage{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Unscoped().Where("seo_product_id = ? ", id).Delete(&model.ProductSeo{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Unscoped().Where("size_product_id = ? ", id).Delete(&model.ProductSize{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (p *ProductRepository) UpdateProduct(product *model.Product) error {
	return p.mysqlDB.Model(product).Update(&product).Error
}

func (p *ProductRepository) FindProductById(id int64) (product *model.Product, err error) {
	return product, p.mysqlDB.Preload("ProductImage").Preload("ProductSize").Preload("ProductSeo").First(product, id).Error
}

func (p *ProductRepository) FindAllProduct() (productAll []model.Product, err error) {
	return productAll, p.mysqlDB.Preload("ProductImage").Preload("ProductSize").Preload("ProductSeo").Find(&productAll).Error
}

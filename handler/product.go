package handler

import (
	"context"

	"github.com/kiven-man/product/domain/model"

	"github.com/kiven-man/product/common"
	"github.com/kiven-man/product/domain/service"
	product "github.com/kiven-man/product/proto"
)

type Product struct {
	ProductDataService service.IProductDataService
}

func (p *Product) AddProduct(ctx context.Context, productInfo *product.ProductInfo, responseProduct *product.ResponseProduct) error {
	product := &model.Product{}
	if err := common.SwapTo(productInfo, product); err != nil {
		return err
	}
	id, err := p.ProductDataService.AddProduct(product)
	if err != nil {
		return err
	}
	responseProduct.ProductId = id
	return nil

}
func (p *Product) DeleteProduct(ctx context.Context, requestID *product.RequestID, response *product.Response) error {
	if err := p.ProductDataService.DeleteProduct(requestID.ProductId); err != nil {
		return err
	}
	response.Msg = "删除成功！"
	return nil
}
func (p *Product) UpdateProduct(ctx context.Context, productInfo *product.ProductInfo, response *product.Response) error {
	product := &model.Product{}
	if err := common.SwapTo(productInfo, product); err != nil {
		return err
	}
	if err := p.ProductDataService.UpdateProduct(product); err != nil {
		return err
	}
	response.Msg = "更新成功！"
	return nil
}
func (p *Product) FindProductById(ctx context.Context, requestID *product.RequestID, productInfo *product.ProductInfo) error {
	product, err := p.ProductDataService.FindProductById(requestID.ProductId)
	if err != nil {
		return err
	}
	if err := common.SwapTo(product, productInfo); err != nil {
		return err
	}
	return nil
}
func (p *Product) FindAllProduct(ctx context.Context, requestAll *product.RequestAll, allProduct *product.AllProduct) error {
	productArr, err := p.ProductDataService.FindAllProduct()
	if err != nil {
		return err
	}
	for _, v := range productArr {
		productInfo := &product.ProductInfo{}
		common.SwapTo(v, productInfo)
		allProduct.ProductInfo = append(allProduct.ProductInfo, productInfo)
	}
	return nil
}

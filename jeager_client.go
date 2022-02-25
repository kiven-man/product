package main

import (
	"context"
	"fmt"

	"github.com/kiven-man/product/common"
	product2 "github.com/kiven-man/product/proto"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

func main() {
	// consul 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.product.client", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 设置服务
	srv := micro.NewService(
		micro.Name("go.micro.service.product.client"),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8085"),
		//	添加注册中心
		micro.Registry(consulRegistry),
		//	绑定链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
	)
	productSrv := product2.NewProductService("go.micro.service.product", srv.Client())
	productAdd := &product2.ProductInfo{
		ProductName:        "imooc",
		ProductSku:         "kevin",
		ProductPrice:       1.1,
		ProductDescription: "imooc-kevin",
		ProductImage: []*product2.ProductImage{
			{
				ImageName: "kevin-image_001",
				ImageCode: "kevin_img_001",
				ImageUrl:  "kevin.com/001",
			}, {
				ImageName: "kevin-image_002",
				ImageCode: "kevin_img_002",
				ImageUrl:  "kevin.com/002",
			},
		},
		ProductSize: []*product2.ProductSize{
			{
				SizeName: "kevin-size",
				SizeCode: "kevin-size-code",
			},
		},
		ProductSeo: &product2.ProductSeo{
			SeoTitle:       "Title_kevin",
			SeoCode:        "code_kevin",
			SeoDescription: "description_kevin",
			SeoKeywords:    "keywords_kevin",
		},
	}
	response, err := productSrv.AddProduct(context.TODO(), productAdd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
}

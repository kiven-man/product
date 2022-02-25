package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kiven-man/product/common"
	"github.com/kiven-man/product/domain/repository"
	service2 "github.com/kiven-man/product/domain/service"
	"github.com/kiven-man/product/handler"
	product2 "github.com/kiven-man/product/proto"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

func main() {
	// consul 配置中心
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	// consul 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.product", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)
	// 数据库设置
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	// 禁止副表
	db.SingularTable(true)
	// 初始化表
	//repository.NewProductRepository(db).InitTable()
	productDataService := service2.NewProductDataService(repository.NewProductRepository(db))

	// 设置服务
	srv := micro.NewService(
		micro.Name("go.micro.service.product"),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8085"),
		//	添加注册中心
		micro.Registry(consulRegistry),
		//	绑定链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	// 初始化服务
	srv.Init()
	// Register handler
	product2.RegisterProductHandler(srv.Server(), &handler.Product{ProductDataService: productDataService})
	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

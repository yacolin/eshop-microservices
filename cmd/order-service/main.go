package main

import (
	"log"
	"os"

	"eshop-microservices/docs"
	_ "eshop-microservices/internal/order-service/api/handlers"
	"eshop-microservices/internal/order-service/app"
)

func main() {
	// 初始化swagger文档
	docs.SwaggerInfo.Title = "订单服务 API"
	docs.SwaggerInfo.Description = "订单服务的API文档"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8081"
	docs.SwaggerInfo.BasePath = "/order/api"

	app, err := app.New(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatalf("app: %v", err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}

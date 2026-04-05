package main

import (
	"log"
	"os"

	"eshop-microservices/docs"
	"eshop-microservices/internal/user-service/app"
)

func main() {
	// 初始化swagger文档
	docs.SwaggerInfo.Title = "用户服务 API"
	docs.SwaggerInfo.Description = "用户服务的API文档"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8082"
	docs.SwaggerInfo.BasePath = "/user/api"

	app, err := app.New(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatalf("app: %v", err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}

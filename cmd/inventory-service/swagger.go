package main

import (
	_ "eshop-microservices/internal/inventory-service/api/handlers"
)

// @title 库存服务 API
// @version 1.0
// @description 库存服务的API文档
// @host localhost:8081
// @BasePath /api

func init() {
	// 导入处理器包，让swag工具能够扫描到它们
}

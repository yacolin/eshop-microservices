package main

import (
	_ "eshop-microservices/internal/user-service/api/handlers"
)

// @title 用户服务 API
// @version 1.0
// @description 用户服务的API文档
// @host localhost:8082
// @BasePath /user/api

func init() {
	// 导入处理器包，让swag工具能够扫描到它们
}

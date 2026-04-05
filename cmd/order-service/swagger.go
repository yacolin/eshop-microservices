package main

import (
	_ "eshop-microservices/internal/order-service/api/handlers"
)

// @title 订单服务 API
// @version 1.0
// @description 订单服务的API文档，支持订单创建、查询、状态更新和取消等功能
// @host localhost:8081
// @BasePath /order/api

func init() {
	// 导入处理器包，让swag工具能够扫描到它们
}

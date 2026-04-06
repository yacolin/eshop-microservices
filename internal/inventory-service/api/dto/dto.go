package dto

import (
	"eshop-microservices/internal/inventory-service/domain/models"
	pkgQuery "eshop-microservices/pkg/query"
)

// ProductListQuery 产品列表查询参数
type ProductListQuery struct {
	pkgQuery.Pagination
	Name   string `form:"name"`              // 产品名称模糊搜索
	SKU    string `form:"sku"`               // SKU精确搜索
	SortBy string `form:"sort_by"`           // 排序字段，例如 name, price, created_at
	Order  string `form:"order,default=asc"` // asc or desc
}

type ProductListResult struct {
	Total int64            `json:"total"`
	List  []models.Product `json:"list"`
}

// InventoryListQuery 库存列表查询参数
type InventoryListQuery struct {
	pkgQuery.Pagination
	ProductName string `form:"product_name"`      // 产品名称模糊搜索
	SKU         string `form:"sku"`               // SKU精确搜索
	Status      string `form:"status"`            // 库存状态
	LowStock    *bool  `form:"low_stock"`         // 是否低库存
	SortBy      string `form:"sort_by"`           // 排序字段，例如 quantity, reserved, created_at
	Order       string `form:"order,default=asc"` // asc or desc
}

type InventoryListResult struct {
	Total int64              `json:"total"`
	List  []models.Inventory `json:"list"`
}

// CreateProductDTO 创建产品请求
type CreateProductDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       int64  `json:"price" binding:"required,min=0"` // 价格，单位：分
	SKU         string `json:"sku" binding:"required"`
}

// UpdateProductDTO 更新产品请求
type UpdateProductDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Price       *int64  `json:"price"`
}

// CreateInventoryDTO 创建库存请求
type CreateInventoryDTO struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=0"`
	Threshold int    `json:"threshold" binding:"min=0"`
}

// UpdateInventoryDTO 更新库存请求
type UpdateInventoryDTO struct {
	Quantity  *int `json:"quantity"`
	Threshold *int `json:"threshold"`
	Reserved  *int `json:"reserved"`
}

// ReserveInventoryDTO 预订库存请求
type ReserveInventoryDTO struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// ReleaseInventoryDTO 释放库存请求
type ReleaseInventoryDTO struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// AdjustInventoryDTO 调整库存请求
type AdjustInventoryDTO struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"` // 正数增加，负数减少
}

// CategoryListQuery 分类列表查询参数
type CategoryListQuery struct {
	pkgQuery.Pagination
	Name     string  `form:"name"`              // 分类名称模糊搜索
	ParentID *string `form:"parent_id"`         // 父分类ID
	SortBy   string  `form:"sort_by"`           // 排序字段，例如 name, created_at
	Order    string  `form:"order,default=asc"` // asc or desc
}

type CategoryListResult struct {
	Total int64             `json:"total"`
	List  []models.Category `json:"list"`
}

// CreateCategoryDTO 创建分类请求
type CreateCategoryDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	ParentID    *string `json:"parent_id"` // 父分类ID，支持层级结构
}

// UpdateCategoryDTO 更新分类请求
type UpdateCategoryDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ParentID    *string `json:"parent_id"`
}

// CreateCommentDTO 创建评论请求
type CreateCommentDTO struct {
	ProductID string  `json:"product_id" binding:"required"`
	Content   string  `json:"content" binding:"required"`
	Rating    int     `json:"rating" binding:"required,min=1,max=5"` // 评分：1-5星
	ParentID  *string `json:"parent_id"` // 父评论ID，用于回复
}

// CommentListQuery 评论列表查询参数
type CommentListQuery struct {
	pkgQuery.Pagination
	ProductID string `form:"product_id" binding:"required"` // 商品ID
	Rating    int    `form:"rating"`                      // 评分筛选
	SortBy    string `form:"sort_by"`                    // 排序字段，例如 rating, created_at
	Order     string `form:"order,default=desc"`         // asc or desc
}

// CommentListResult 评论列表结果
type CommentListResult struct {
	Total    int64            `json:"total"`
	List     []models.Comment `json:"list"`
	AvgRating float64          `json:"avg_rating"` // 商品平均评分
}

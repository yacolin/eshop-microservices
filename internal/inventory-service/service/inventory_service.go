package service

import (
	"context"
	"fmt"

	"eshop-microservices/internal/inventory-service/api/dto"
	"eshop-microservices/internal/inventory-service/domain/models"
	"eshop-microservices/internal/inventory-service/domain/repositories"
	"eshop-microservices/pkg/errcode"
	"eshop-microservices/pkg/query"
)

// InventoryService 库存业务
type InventoryService struct {
	repo repositories.InventoryRepository
}

// CreateCategory 创建分类
func (s *InventoryService) CreateCategory(ctx context.Context, req dto.CreateCategoryDTO) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	if err := s.repo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoryByID 获取分类详情
func (s *InventoryService) GetCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return category, nil
}

// UpdateCategory 更新分类
func (s *InventoryService) UpdateCategory(ctx context.Context, id string, req dto.UpdateCategoryDTO) (*models.Category, error) {
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}

	if err := s.repo.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory 删除分类
func (s *InventoryService) DeleteCategory(ctx context.Context, id string) error {
	return s.repo.DeleteCategory(ctx, id)
}

// ListCategories 分类列表
func (s *InventoryService) ListCategories(ctx context.Context, q dto.CategoryListQuery) (*query.ListResult[models.Category], error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 {
		q.Size = 10
	}
	if q.Size > 100 {
		q.Size = 100
	}
	offset := (q.Page - 1) * q.Size

	list, err := s.repo.ListCategories(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountCategories(ctx, q)
	if err != nil {
		return nil, err
	}

	return &query.ListResult[models.Category]{
		List:  list,
		Total: total,
	}, nil
}

// NewInventoryService 创建库存服务
func NewInventoryService(repo repositories.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

// CreateProduct 创建产品
func (s *InventoryService) CreateProduct(ctx context.Context, req dto.CreateProductDTO) (*models.Product, error) {
	// 检查SKU是否已存在
	existingProduct, _ := s.repo.GetProductBySKU(ctx, req.SKU)
	if existingProduct != nil {
		return nil, fmt.Errorf("product with SKU %s already exists", req.SKU)
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.SKU,
	}

	if err := s.repo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID 获取产品详情
func (s *InventoryService) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return product, nil
}

// UpdateProduct 更新产品
func (s *InventoryService) UpdateProduct(ctx context.Context, id string, req dto.UpdateProductDTO) (*models.Product, error) {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}

	if err := s.repo.UpdateProduct(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct 删除产品
func (s *InventoryService) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.DeleteProduct(ctx, id)
}

// ListProducts 产品列表
func (s *InventoryService) ListProducts(ctx context.Context, q dto.ProductListQuery) (*query.ListResult[models.Product], error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 {
		q.Size = 10
	}
	if q.Size > 100 {
		q.Size = 100
	}
	offset := (q.Page - 1) * q.Size

	list, err := s.repo.ListProducts(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountProducts(ctx, q)
	if err != nil {
		return nil, err
	}

	return &query.ListResult[models.Product]{
		List:  list,
		Total: total,
	}, nil
}

// CreateInventory 创建库存记录
func (s *InventoryService) CreateInventory(ctx context.Context, req dto.CreateInventoryDTO) (*models.Inventory, error) {
	// 检查产品是否存在
	_, err := s.repo.GetProductByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product with ID %s does not exist", req.ProductID)
	}

	// 检查库存记录是否已存在
	existingInventory, _ := s.repo.GetInventoryByProductID(ctx, req.ProductID)
	if existingInventory != nil {
		return nil, fmt.Errorf("inventory record for product %s already exists", req.ProductID)
	}

	inventory := &models.Inventory{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Threshold: req.Threshold,
	}
	inventory.UpdateStatus() // 设置初始状态

	if err := s.repo.CreateInventory(ctx, inventory); err != nil {
		return nil, err
	}

	return inventory, nil
}

// GetInventoryByID 获取库存详情
func (s *InventoryService) GetInventoryByID(ctx context.Context, id string) (*models.Inventory, error) {
	inventory, err := s.repo.GetInventoryByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return inventory, nil
}

// GetInventoryByProductID 根据产品ID获取库存
func (s *InventoryService) GetInventoryByProductID(ctx context.Context, productID string) (*models.Inventory, error) {
	inventory, err := s.repo.GetInventoryByProductID(ctx, productID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return inventory, nil
}

// UpdateInventory 更新库存
func (s *InventoryService) UpdateInventory(ctx context.Context, id string, req dto.UpdateInventoryDTO) (*models.Inventory, error) {
	inventory, err := s.repo.GetInventoryByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}

	if req.Quantity != nil {
		inventory.Quantity = *req.Quantity
	}
	if req.Threshold != nil {
		inventory.Threshold = *req.Threshold
	}
	if req.Reserved != nil {
		inventory.Reserved = *req.Reserved
	}

	if err := s.repo.UpdateInventory(ctx, inventory); err != nil {
		return nil, err
	}

	return inventory, nil
}

// DeleteInventory 删除库存
func (s *InventoryService) DeleteInventory(ctx context.Context, id string) error {
	return s.repo.DeleteInventory(ctx, id)
}

// ListInventories 库存列表
func (s *InventoryService) ListInventories(ctx context.Context, q dto.InventoryListQuery) (*query.ListResult[models.Inventory], error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 {
		q.Size = 10
	}
	if q.Size > 100 {
		q.Size = 100
	}
	offset := (q.Page - 1) * q.Size

	list, err := s.repo.ListInventories(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountInventories(ctx, q)
	if err != nil {
		return nil, err
	}

	return &query.ListResult[models.Inventory]{
		List:  list,
		Total: total,
	}, nil
}

// ReserveInventory 预订库存
func (s *InventoryService) ReserveInventory(ctx context.Context, req dto.ReserveInventoryDTO) error {
	inventory, err := s.repo.GetInventoryByProductID(ctx, req.ProductID)
	if err != nil {
		return fmt.Errorf("inventory not found for product %s", req.ProductID)
	}

	if inventory.Quantity-inventory.Reserved < req.Quantity {
		return fmt.Errorf("insufficient stock for reservation: requested %d, available %d",
			req.Quantity, inventory.Quantity-inventory.Reserved)
	}

	inventory.Reserved += req.Quantity
	inventory.UpdateStatus()

	if err := s.repo.UpdateInventory(ctx, inventory); err != nil {
		return err
	}

	return nil
}

// ReleaseInventory 释放库存
func (s *InventoryService) ReleaseInventory(ctx context.Context, req dto.ReleaseInventoryDTO) error {
	inventory, err := s.repo.GetInventoryByProductID(ctx, req.ProductID)
	if err != nil {
		return fmt.Errorf("inventory not found for product %s", req.ProductID)
	}

	if inventory.Reserved < req.Quantity {
		return fmt.Errorf("not enough reserved stock to release: requested %d, reserved %d",
			req.Quantity, inventory.Reserved)
	}

	inventory.Reserved -= req.Quantity
	inventory.UpdateStatus()

	if err := s.repo.UpdateInventory(ctx, inventory); err != nil {
		return err
	}

	return nil
}

// AdjustInventory 调整库存数量
func (s *InventoryService) AdjustInventory(ctx context.Context, req dto.AdjustInventoryDTO) error {
	if req.Quantity == 0 {
		return fmt.Errorf("adjustment quantity cannot be zero")
	}

	inventory, err := s.repo.GetInventoryByProductID(ctx, req.ProductID)
	if err != nil {
		return fmt.Errorf("inventory not found for product %s", req.ProductID)
	}

	newQuantity := inventory.Quantity + req.Quantity
	if newQuantity < 0 {
		return fmt.Errorf("adjustment would result in negative inventory: current %d, adjustment %d",
			inventory.Quantity, req.Quantity)
	}

	inventory.Quantity = newQuantity
	inventory.UpdateStatus()

	if err := s.repo.UpdateInventory(ctx, inventory); err != nil {
		return err
	}

	return nil
}

// CheckAvailability 检查产品是否有足够库存
func (s *InventoryService) CheckAvailability(ctx context.Context, productID string, quantity int) (bool, error) {
	inventory, err := s.repo.GetInventoryByProductID(ctx, productID)
	if err != nil {
		return false, fmt.Errorf("inventory not found for product %s", productID)
	}

	available := inventory.Quantity - inventory.Reserved
	return available >= quantity, nil
}

// CreateComment 创建评论
func (s *InventoryService) CreateComment(ctx context.Context, req dto.CreateCommentDTO, userID string) (*models.Comment, error) {
	// 检查商品是否存在
	_, err := s.repo.GetProductByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product with ID %s does not exist", req.ProductID)
	}

	// 如果是回复评论，检查父评论是否存在
	if req.ParentID != nil {
		_, err := s.repo.GetCommentByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent comment with ID %s does not exist", *req.ParentID)
		}
	}

	comment := &models.Comment{
		ProductID: req.ProductID,
		UserID:    userID,
		Content:   req.Content,
		Rating:    req.Rating,
		ParentID:  req.ParentID,
	}

	if err := s.repo.CreateComment(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

// ListComments 获取评论列表
func (s *InventoryService) ListComments(ctx context.Context, q dto.CommentListQuery) (*dto.CommentListResult, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 {
		q.Size = 10
	}
	if q.Size > 100 {
		q.Size = 100
	}
	offset := (q.Page - 1) * q.Size

	list, err := s.repo.ListComments(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountComments(ctx, q)
	if err != nil {
		return nil, err
	}

	// 计算平均评分
	avgRating, err := s.repo.GetAverageRating(ctx, q.ProductID)
	if err != nil {
		avgRating = 0
	}

	return &dto.CommentListResult{
		Total:     total,
		List:      list,
		AvgRating: avgRating,
	}, nil
}

// DeleteComment 删除评论
func (s *InventoryService) DeleteComment(ctx context.Context, id string) error {
	// 检查评论是否存在
	_, err := s.repo.GetCommentByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}

	return s.repo.DeleteComment(ctx, id)
}

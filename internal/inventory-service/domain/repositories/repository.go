package repositories

import (
	"context"

	"eshop-microservices/internal/inventory-service/api/dto"
	"eshop-microservices/internal/inventory-service/domain/models"

	"gorm.io/gorm"
)

// Repository 组合仓储接口，包含所有子仓储接口
type Repository interface {
	ProductRepository
	InventoryRepositoryInterface
	CategoryRepository
}

// repository 组合仓储实现
type repository struct {
	productRepo  ProductRepository
	inventoryRepo InventoryRepositoryInterface
	categoryRepo CategoryRepository
}

// NewRepository 创建组合仓储
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		productRepo:  NewProductRepository(db),
		inventoryRepo: NewInventoryRepositoryImpl(db),
		categoryRepo: NewCategoryRepository(db),
	}
}

// Product Operations
func (r *repository) CreateProduct(ctx context.Context, product *models.Product) error {
	return r.productRepo.CreateProduct(ctx, product)
}

func (r *repository) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	return r.productRepo.GetProductByID(ctx, id)
}

func (r *repository) GetProductBySKU(ctx context.Context, sku string) (*models.Product, error) {
	return r.productRepo.GetProductBySKU(ctx, sku)
}

func (r *repository) UpdateProduct(ctx context.Context, product *models.Product) error {
	return r.productRepo.UpdateProduct(ctx, product)
}

func (r *repository) DeleteProduct(ctx context.Context, id string) error {
	return r.productRepo.DeleteProduct(ctx, id)
}

func (r *repository) ListProducts(ctx context.Context, q dto.ProductListQuery, offset, limit int) ([]models.Product, error) {
	return r.productRepo.ListProducts(ctx, q, offset, limit)
}

func (r *repository) CountProducts(ctx context.Context, q dto.ProductListQuery) (int64, error) {
	return r.productRepo.CountProducts(ctx, q)
}

// Inventory Operations
func (r *repository) CreateInventory(ctx context.Context, inventory *models.Inventory) error {
	return r.inventoryRepo.CreateInventory(ctx, inventory)
}

func (r *repository) GetInventoryByID(ctx context.Context, id string) (*models.Inventory, error) {
	return r.inventoryRepo.GetInventoryByID(ctx, id)
}

func (r *repository) GetInventoryByProductID(ctx context.Context, productID string) (*models.Inventory, error) {
	return r.inventoryRepo.GetInventoryByProductID(ctx, productID)
}

func (r *repository) UpdateInventory(ctx context.Context, inventory *models.Inventory) error {
	return r.inventoryRepo.UpdateInventory(ctx, inventory)
}

func (r *repository) UpdateInventoryQuantity(ctx context.Context, productID string, quantityChange int) error {
	return r.inventoryRepo.UpdateInventoryQuantity(ctx, productID, quantityChange)
}

func (r *repository) DeleteInventory(ctx context.Context, id string) error {
	return r.inventoryRepo.DeleteInventory(ctx, id)
}

func (r *repository) ListInventories(ctx context.Context, q dto.InventoryListQuery, offset, limit int) ([]models.Inventory, error) {
	return r.inventoryRepo.ListInventories(ctx, q, offset, limit)
}

func (r *repository) CountInventories(ctx context.Context, q dto.InventoryListQuery) (int64, error) {
	return r.inventoryRepo.CountInventories(ctx, q)
}

// Category Operations
func (r *repository) CreateCategory(ctx context.Context, category *models.Category) error {
	return r.categoryRepo.CreateCategory(ctx, category)
}

func (r *repository) GetCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	return r.categoryRepo.GetCategoryByID(ctx, id)
}

func (r *repository) UpdateCategory(ctx context.Context, category *models.Category) error {
	return r.categoryRepo.UpdateCategory(ctx, category)
}

func (r *repository) DeleteCategory(ctx context.Context, id string) error {
	return r.categoryRepo.DeleteCategory(ctx, id)
}

func (r *repository) ListCategories(ctx context.Context, q dto.CategoryListQuery, offset, limit int) ([]models.Category, error) {
	return r.categoryRepo.ListCategories(ctx, q, offset, limit)
}

func (r *repository) CountCategories(ctx context.Context, q dto.CategoryListQuery) (int64, error) {
	return r.categoryRepo.CountCategories(ctx, q)
}
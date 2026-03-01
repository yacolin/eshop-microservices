package repositories

import (
	"context"

	"eshop-microservices/internal/inventory-service/api/dto"
	"eshop-microservices/internal/inventory-service/domain/models"

	"gorm.io/gorm"
)

// InventoryRepository 库存仓储接口
type InventoryRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, q dto.ProductListQuery, offset, limit int) ([]models.Product, error)
	CountProducts(ctx context.Context, q dto.ProductListQuery) (int64, error)

	CreateInventory(ctx context.Context, inventory *models.Inventory) error
	GetInventoryByID(ctx context.Context, id string) (*models.Inventory, error)
	GetInventoryByProductID(ctx context.Context, productID string) (*models.Inventory, error)
	UpdateInventory(ctx context.Context, inventory *models.Inventory) error
	UpdateInventoryQuantity(ctx context.Context, productID string, quantityChange int) error
	DeleteInventory(ctx context.Context, id string) error
	ListInventories(ctx context.Context, q dto.InventoryListQuery, offset, limit int) ([]models.Inventory, error)
	CountInventories(ctx context.Context, q dto.InventoryListQuery) (int64, error)

	CreateCategory(ctx context.Context, category *models.Category) error
	GetCategoryByID(ctx context.Context, id string) (*models.Category, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	DeleteCategory(ctx context.Context, id string) error
	ListCategories(ctx context.Context, q dto.CategoryListQuery, offset, limit int) ([]models.Category, error)
	CountCategories(ctx context.Context, q dto.CategoryListQuery) (int64, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository 创建库存仓储
func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

// Product Operations
func (r *inventoryRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *inventoryRepository) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *inventoryRepository) GetProductBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *inventoryRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *inventoryRepository) DeleteProduct(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, "id = ?", id).Error
}

func (r *inventoryRepository) ListProducts(ctx context.Context, q dto.ProductListQuery, offset, limit int) ([]models.Product, error) {
	var products []models.Product
	db := r.ApplyProductQuery(ctx, q)
	err := db.Offset(offset).Limit(limit).Find(&products).Error
	return products, err
}

func (r *inventoryRepository) CountProducts(ctx context.Context, q dto.ProductListQuery) (int64, error) {
	var count int64
	db := r.ApplyProductQuery(ctx, q)
	err := db.Count(&count).Error
	return count, err
}

func (r *inventoryRepository) ApplyProductQuery(ctx context.Context, q dto.ProductListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.Product{})
	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.SKU != "" {
		db = db.Where("sku = ?", q.SKU)
	}

	order := "id asc"
	if q.SortBy != "" {
		ord := q.Order
		if ord != "asc" && ord != "desc" {
			ord = "asc"
		}
		order = q.SortBy + " " + ord
	}

	db = db.Order(order)
	return db
}

// Inventory Operations
func (r *inventoryRepository) CreateInventory(ctx context.Context, inventory *models.Inventory) error {
	return r.db.WithContext(ctx).Create(inventory).Error
}

func (r *inventoryRepository) GetInventoryByID(ctx context.Context, id string) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.WithContext(ctx).Preload("Product").Where("id = ?", id).First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) GetInventoryByProductID(ctx context.Context, productID string) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.WithContext(ctx).Preload("Product").Where("product_id = ?", productID).First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) UpdateInventory(ctx context.Context, inventory *models.Inventory) error {
	inventory.UpdateStatus() // 更新状态
	return r.db.WithContext(ctx).Save(inventory).Error
}

func (r *inventoryRepository) UpdateInventoryQuantity(ctx context.Context, productID string, quantityChange int) error {
	tx := r.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	inventory, err := r.GetInventoryByProductID(ctx, productID)
	if err != nil {
		return err
	}

	newQuantity := inventory.Quantity + quantityChange
	if newQuantity < 0 {
		return gorm.ErrRecordNotFound // 或者自定义错误表示库存不足
	}

	inventory.Quantity = newQuantity
	inventory.UpdateStatus() // 更新状态

	if err := tx.Save(inventory).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

func (r *inventoryRepository) DeleteInventory(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Inventory{}, "id = ?", id).Error
}

func (r *inventoryRepository) ListInventories(ctx context.Context, q dto.InventoryListQuery, offset, limit int) ([]models.Inventory, error) {
	var inventories []models.Inventory
	db := r.ApplyInventoryQuery(ctx, q)
	err := db.Preload("Product").Offset(offset).Limit(limit).Find(&inventories).Error
	return inventories, err
}

func (r *inventoryRepository) CountInventories(ctx context.Context, q dto.InventoryListQuery) (int64, error) {
	var count int64
	db := r.ApplyInventoryQuery(ctx, q)
	err := db.Count(&count).Error
	return count, err
}

func (r *inventoryRepository) ApplyInventoryQuery(ctx context.Context, q dto.InventoryListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.Inventory{})
	db = db.Joins("JOIN products ON inventories.product_id = products.id")

	if q.ProductName != "" {
		db = db.Where("products.name LIKE ?", "%"+q.ProductName+"%")
	}
	if q.SKU != "" {
		db = db.Where("products.sku = ?", q.SKU)
	}
	if q.Status != "" {
		db = db.Where("inventories.status = ?", q.Status)
	}
	if q.LowStock != nil && *q.LowStock {
		db = db.Where("inventories.quantity <= inventories.threshold AND inventories.quantity > 0")
	}

	order := "inventories.id asc"
	if q.SortBy != "" {
		ord := q.Order
		if ord != "asc" && ord != "desc" {
			ord = "asc"
		}
		// Map sort field to appropriate table column
		switch q.SortBy {
		case "name":
			order = "products.name " + ord
		case "sku":
			order = "products.sku " + ord
		default:
			order = "inventories." + q.SortBy + " " + ord
		}
	}

	db = db.Order(order)
	return db
}

// Category Operations
func (r *inventoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *inventoryRepository) GetCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).Preload("Parent").Preload("Children").Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *inventoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *inventoryRepository) DeleteCategory(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Category{}, "id = ?", id).Error
}

func (r *inventoryRepository) ListCategories(ctx context.Context, q dto.CategoryListQuery, offset, limit int) ([]models.Category, error) {
	var categories []models.Category
	db := r.ApplyCategoryQuery(ctx, q)
	// 预加载父分类和子分类
	err := db.Preload("Parent").Preload("Children").Offset(offset).Limit(limit).Find(&categories).Error
	return categories, err
}

func (r *inventoryRepository) CountCategories(ctx context.Context, q dto.CategoryListQuery) (int64, error) {
	var count int64
	db := r.ApplyCategoryQuery(ctx, q)
	err := db.Count(&count).Error
	return count, err
}

func (r *inventoryRepository) ApplyCategoryQuery(ctx context.Context, q dto.CategoryListQuery) *gorm.DB {
	// 使用左连接来处理可能的空值
	db := r.db.WithContext(ctx).Model(&models.Category{})
	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.ParentID != nil {
		db = db.Where("parent_id = ?", *q.ParentID)
	}

	order := "id asc"
	if q.SortBy != "" {
		ord := q.Order
		if ord != "asc" && ord != "desc" {
			ord = "asc"
		}
		order = q.SortBy + " " + ord
	}

	db = db.Order(order)
	return db
}

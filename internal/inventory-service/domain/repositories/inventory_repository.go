package repositories

import (
	"context"

	"eshop-microservices/internal/inventory-service/api/dto"
	"eshop-microservices/internal/inventory-service/domain/models"

	"gorm.io/gorm"
)

// InventoryRepository 库存仓储接口（保持向后兼容）
type InventoryRepository interface {
	ProductRepository
	InventoryRepositoryInterface
	CategoryRepository
}

// InventoryRepositoryInterface 库存仓储核心接口
type InventoryRepositoryInterface interface {
	CreateInventory(ctx context.Context, inventory *models.Inventory) error
	GetInventoryByID(ctx context.Context, id string) (*models.Inventory, error)
	GetInventoryByProductID(ctx context.Context, productID string) (*models.Inventory, error)
	UpdateInventory(ctx context.Context, inventory *models.Inventory) error
	UpdateInventoryQuantity(ctx context.Context, productID string, quantityChange int) error
	DeleteInventory(ctx context.Context, id string) error
	ListInventories(ctx context.Context, q dto.InventoryListQuery, offset, limit int) ([]models.Inventory, error)
	CountInventories(ctx context.Context, q dto.InventoryListQuery) (int64, error)
}

type inventoryRepository struct {
	db *gorm.DB
}


// NewInventoryRepositoryImpl 创建库存仓储实现
func NewInventoryRepositoryImpl(db *gorm.DB) InventoryRepositoryInterface {
	return &inventoryRepository{db: db}
}

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

// NewInventoryRepository 创建库存仓储（保持向后兼容）
func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return NewRepository(db)
}
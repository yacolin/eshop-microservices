package repositories

import (
	"context"

	"eshop-microservices/internal/inventory-service/api/dto"
	"eshop-microservices/internal/inventory-service/domain/models"

	"gorm.io/gorm"
)

// ProductRepository 产品仓储接口
type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, q dto.ProductListQuery, offset, limit int) ([]models.Product, error)
	CountProducts(ctx context.Context, q dto.ProductListQuery) (int64, error)
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建产品仓储
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetProductBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) DeleteProduct(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, "id = ?", id).Error
}

func (r *productRepository) ListProducts(ctx context.Context, q dto.ProductListQuery, offset, limit int) ([]models.Product, error) {
	var products []models.Product
	db := r.ApplyProductQuery(ctx, q)

	// 预加载分类Category 和多对多关系Categories
	err := db.Preload("Category").Preload("Categories").Offset(offset).Limit(limit).Find(&products).Error
	return products, err
}

func (r *productRepository) CountProducts(ctx context.Context, q dto.ProductListQuery) (int64, error) {
	var count int64
	db := r.ApplyProductQuery(ctx, q)
	err := db.Count(&count).Error
	return count, err
}

func (r *productRepository) ApplyProductQuery(ctx context.Context, q dto.ProductListQuery) *gorm.DB {
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

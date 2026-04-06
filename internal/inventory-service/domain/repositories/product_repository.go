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

	// 评论相关方法
	CreateComment(ctx context.Context, comment *models.Comment) error
	GetCommentByID(ctx context.Context, id string) (*models.Comment, error)
	ListComments(ctx context.Context, q dto.CommentListQuery, offset, limit int) ([]models.Comment, error)
	CountComments(ctx context.Context, q dto.CommentListQuery) (int64, error)
	GetAverageRating(ctx context.Context, productID string) (float64, error)
	DeleteComment(ctx context.Context, id string) error
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

// CreateComment 创建评论
func (r *productRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// GetCommentByID 获取评论详情
func (r *productRepository) GetCommentByID(ctx context.Context, id string) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListComments 获取评论列表
func (r *productRepository) ListComments(ctx context.Context, q dto.CommentListQuery, offset, limit int) ([]models.Comment, error) {
	var comments []models.Comment
	db := r.ApplyCommentQuery(ctx, q)
	err := db.Offset(offset).Limit(limit).Find(&comments).Error
	return comments, err
}

// CountComments 统计评论数量
func (r *productRepository) CountComments(ctx context.Context, q dto.CommentListQuery) (int64, error) {
	var count int64
	db := r.ApplyCommentQuery(ctx, q)
	err := db.Count(&count).Error
	return count, err
}

// GetAverageRating 计算平均评分
func (r *productRepository) GetAverageRating(ctx context.Context, productID string) (float64, error) {
	type Result struct {
		AvgRating float64
	}
	var result Result
	err := r.db.WithContext(ctx).Model(&models.Comment{}).
		Select("COALESCE(AVG(rating), 0) as avg_rating").
		Where("product_id = ?", productID).
		Scan(&result).Error
	return result.AvgRating, err
}

// ApplyCommentQuery 应用评论查询条件
func (r *productRepository) ApplyCommentQuery(ctx context.Context, q dto.CommentListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.Comment{})
	db = db.Where("product_id = ?", q.ProductID)

	if q.Rating > 0 {
		db = db.Where("rating = ?", q.Rating)
	}

	order := "created_at desc"
	if q.SortBy != "" {
		ord := q.Order
		if ord != "asc" && ord != "desc" {
			ord = "desc"
		}
		order = q.SortBy + " " + ord
	}

	db = db.Order(order)
	return db
}

// DeleteComment 删除评论
func (r *productRepository) DeleteComment(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Comment{}, "id = ?", id).Error
}

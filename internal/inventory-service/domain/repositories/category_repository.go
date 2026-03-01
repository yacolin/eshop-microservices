package repositories

import (
	"context"

	"eshop-microservices/internal/inventory-service/api/dto"
	"eshop-microservices/internal/inventory-service/domain/models"

	"gorm.io/gorm"
)

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *models.Category) error
	GetCategoryByID(ctx context.Context, id string) (*models.Category, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	DeleteCategory(ctx context.Context, id string) error
	ListCategories(ctx context.Context, q dto.CategoryListQuery, offset, limit int) ([]models.Category, error)
	CountCategories(ctx context.Context, q dto.CategoryListQuery) (int64, error)
}

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).Preload("Parent").Preload("Children").Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Category{}, "id = ?", id).Error
}

func (r *categoryRepository) ListCategories(ctx context.Context, q dto.CategoryListQuery, offset, limit int) ([]models.Category, error) {
	var categories []models.Category
	db := r.ApplyCategoryQuery(ctx, q)
	// 预加载父分类和子分类
	err := db.Preload("Parent").Preload("Children").Offset(offset).Limit(limit).Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) CountCategories(ctx context.Context, q dto.CategoryListQuery) (int64, error) {
	var count int64
	db := r.ApplyCategoryQuery(ctx, q)
	err := db.Count(&count).Error
	return count, err
}

func (r *categoryRepository) ApplyCategoryQuery(ctx context.Context, q dto.CategoryListQuery) *gorm.DB {
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
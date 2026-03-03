package repositories

import (
	"context"

	"eshop-microservices/internal/order-service/api/dto"
	"eshop-microservices/internal/order-service/domain/models"

	"gorm.io/gorm"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id string) (*models.Order, error)
	List(ctx context.Context, customerID string, limit, offset int) ([]models.Order, int64, error)
	UpdateStatus(ctx context.Context, id, status string) error
	Delete(ctx context.Context, id string) error

	ListByQuery(ctx context.Context, q dto.OrderListQuery, offset, limit int) ([]models.Order, error)
	CountByQuery(ctx context.Context, q dto.OrderListQuery) (int64, error)
}

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	// GORM will automatically create related OrderItem records when the
	// parent Order is created (since Items is a has-many relation). The
	// previous implementation explicitly called Create on order.Items a
	// second time, resulting in duplicate primary key errors when the
	// BeforeCreate hook generated IDs on the first insert. We now rely on
	// the single call below to insert both order and items in one shot.
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *orderRepository) GetByID(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) List(ctx context.Context, customerID string, limit, offset int) ([]models.Order, int64, error) {
	var list []models.Order
	query := r.db.WithContext(ctx).Model(&models.Order{})
	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Preload("Items").Limit(limit).Offset(offset).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", id).Delete(&models.OrderItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Order{}, "id = ?", id).Error
	})
}

func (r *orderRepository) ListByQuery(ctx context.Context, q dto.OrderListQuery, offset, limit int) ([]models.Order, error) {
	var list []models.Order
	db := r.ApplyQuery(ctx, q)
	if err := db.Preload("Items").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r *orderRepository) CountByQuery(ctx context.Context, q dto.OrderListQuery) (int64, error) {
	var count int64
	db := r.ApplyQuery(ctx, q)
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *orderRepository) ApplyQuery(ctx context.Context, q dto.OrderListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.Order{})
	if q.CustomerID != nil {
		db = db.Where("customer_id = ?", q.CustomerID)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}

	if q.MinPrice != nil {
		db = db.Where("total_amount >= ?", q.MinPrice)
	}

	if q.MaxPrice != nil {
		db = db.Where("total_amount <= ?", q.MaxPrice)
	}

	order := "id asc"
	if q.SortBy != "" {
		ord := q.Order
		if ord != "asc" && ord != "desc" {
			ord = "asc"
		}
		order = q.SortBy + " " + q.Order
	}

	db = db.Order(order)

	return db
}

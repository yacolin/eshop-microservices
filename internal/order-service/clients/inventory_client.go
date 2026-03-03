package clients

import (
	"context"
	"fmt"
	"time"

	"eshop-microservices/api/proto/inventorypb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InventoryClient 库存服务 gRPC 客户端
type InventoryClient struct {
	conn   *grpc.ClientConn
	client inventorypb.InventoryServiceClient
}

// NewInventoryClient 创建库存服务客户端
func NewInventoryClient(target string) (*InventoryClient, error) {
	// 设置连接选项
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory service: %w", err)
	}

	client := inventorypb.NewInventoryServiceClient(conn)
	return &InventoryClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close 关闭连接
func (c *InventoryClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ReserveStock 预占库存
func (c *InventoryClient) ReserveStock(ctx context.Context, orderID string, items []*inventorypb.StockItem) (*inventorypb.ReserveStockResponse, error) {
	req := &inventorypb.ReserveStockRequest{
		OrderId: orderID,
		Items:   items,
	}

	resp, err := c.client.ReserveStock(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}

	return resp, nil
}

// ConfirmDeduct 确认扣减库存
func (c *InventoryClient) ConfirmDeduct(ctx context.Context, orderID, reservationID string, items []*inventorypb.StockItem) (*inventorypb.ConfirmDeductResponse, error) {
	req := &inventorypb.ConfirmDeductRequest{
		OrderId:       orderID,
		ReservationId: reservationID,
		Items:         items,
	}

	resp, err := c.client.ConfirmDeduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm deduct: %w", err)
	}

	return resp, nil
}

// ReleaseStock 释放库存
func (c *InventoryClient) ReleaseStock(ctx context.Context, orderID, reservationID string, items []*inventorypb.StockItem) (*inventorypb.ReleaseStockResponse, error) {
	req := &inventorypb.ReleaseStockRequest{
		OrderId:       orderID,
		ReservationId: reservationID,
		Items:         items,
	}

	resp, err := c.client.ReleaseStock(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to release stock: %w", err)
	}

	return resp, nil
}

// CheckStockAvailability 检查库存可用性
func (c *InventoryClient) CheckStockAvailability(ctx context.Context, productID string, quantity int32) (*inventorypb.CheckStockResponse, error) {
	req := &inventorypb.CheckStockRequest{
		ProductId: productID,
		Quantity:  quantity,
	}

	resp, err := c.client.CheckStockAvailability(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to check stock availability: %w", err)
	}

	return resp, nil
}

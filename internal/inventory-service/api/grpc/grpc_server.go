package grpc

import (
	"context"
	"fmt"
	"net"

	"eshop-microservices/api/proto/inventorypb"
	"eshop-microservices/internal/inventory-service/domain/repositories"
	"eshop-microservices/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server gRPC 服务器
type Server struct {
	inventorypb.UnimplementedInventoryServiceServer
	grpcServer *grpc.Server
	repo       repositories.InventoryRepository
	port       string
}

// NewServer 创建 gRPC 服务器
func NewServer(repo repositories.InventoryRepository, port string) *Server {
	return &Server{
		repo: repo,
		port: port,
	}
}

// Start 启动 gRPC 服务器
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.grpcServer = grpc.NewServer()
	inventorypb.RegisterInventoryServiceServer(s.grpcServer, s)

	logger.Info("gRPC server starting", zap.String("port", s.port))
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

// Stop 停止 gRPC 服务器
func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
		logger.Info("gRPC server stopped")
	}
}

// ReserveStock 预占库存
func (s *Server) ReserveStock(ctx context.Context, req *inventorypb.ReserveStockRequest) (*inventorypb.ReserveStockResponse, error) {
	logger.Info("ReserveStock called", zap.String("order_id", req.OrderId))

	// 检查所有商品的库存
	for _, item := range req.Items {
		inventory, err := s.repo.GetInventoryByProductID(ctx, item.ProductId)
		if err != nil {
			return &inventorypb.ReserveStockResponse{
				Success: false,
				Message: fmt.Sprintf("inventory not found for product %s", item.ProductId),
			}, nil
		}

		available := inventory.Quantity - inventory.Reserved
		if available < int(item.Quantity) {
			return &inventorypb.ReserveStockResponse{
				Success: false,
				Message: fmt.Sprintf("insufficient stock for product %s: requested %d, available %d",
					item.ProductId, item.Quantity, available),
			}, nil
		}
	}

	// 预占库存
	for _, item := range req.Items {
		inventory, err := s.repo.GetInventoryByProductID(ctx, item.ProductId)
		if err != nil {
			return &inventorypb.ReserveStockResponse{
				Success: false,
				Message: fmt.Sprintf("inventory not found for product %s", item.ProductId),
			}, nil
		}

		inventory.Reserved += int(item.Quantity)
		inventory.UpdateStatus()

		if err := s.repo.UpdateInventory(ctx, inventory); err != nil {
			return &inventorypb.ReserveStockResponse{
				Success: false,
				Message: fmt.Sprintf("failed to reserve stock for product %s: %v", item.ProductId, err),
			}, nil
		}
	}

	return &inventorypb.ReserveStockResponse{
		Success:       true,
		Message:       "stock reserved successfully",
		ReservationId: req.OrderId, // 使用订单ID作为预留ID
	}, nil
}

// ConfirmDeduct 确认扣减库存（最终扣减）
func (s *Server) ConfirmDeduct(ctx context.Context, req *inventorypb.ConfirmDeductRequest) (*inventorypb.ConfirmDeductResponse, error) {
	logger.Info("ConfirmDeduct called", zap.String("order_id", req.OrderId))

	for _, item := range req.Items {
		inventory, err := s.repo.GetInventoryByProductID(ctx, item.ProductId)
		if err != nil {
			return &inventorypb.ConfirmDeductResponse{
				Success: false,
				Message: fmt.Sprintf("inventory not found for product %s", item.ProductId),
			}, nil
		}

		// 检查是否有足够的预占库存
		if inventory.Reserved < int(item.Quantity) {
			return &inventorypb.ConfirmDeductResponse{
				Success: false,
				Message: fmt.Sprintf("not enough reserved stock for product %s: requested %d, reserved %d",
					item.ProductId, item.Quantity, inventory.Reserved),
			}, nil
		}

		// 扣减实际库存和预占库存
		inventory.Quantity -= int(item.Quantity)
		inventory.Reserved -= int(item.Quantity)
		inventory.UpdateStatus()

		if err := s.repo.UpdateInventory(ctx, inventory); err != nil {
			return &inventorypb.ConfirmDeductResponse{
				Success: false,
				Message: fmt.Sprintf("failed to deduct stock for product %s: %v", item.ProductId, err),
			}, nil
		}
	}

	return &inventorypb.ConfirmDeductResponse{
		Success: true,
		Message: "stock deducted successfully",
	}, nil
}

// ReleaseStock 释放预占库存
func (s *Server) ReleaseStock(ctx context.Context, req *inventorypb.ReleaseStockRequest) (*inventorypb.ReleaseStockResponse, error) {
	logger.Info("ReleaseStock called", zap.String("order_id", req.OrderId))

	for _, item := range req.Items {
		inventory, err := s.repo.GetInventoryByProductID(ctx, item.ProductId)
		if err != nil {
			return &inventorypb.ReleaseStockResponse{
				Success: false,
				Message: fmt.Sprintf("inventory not found for product %s", item.ProductId),
			}, nil
		}

		// 检查是否有足够的预占库存可释放
		if inventory.Reserved < int(item.Quantity) {
			return &inventorypb.ReleaseStockResponse{
				Success: false,
				Message: fmt.Sprintf("not enough reserved stock to release for product %s: requested %d, reserved %d",
					item.ProductId, item.Quantity, inventory.Reserved),
			}, nil
		}

		// 释放预占库存
		inventory.Reserved -= int(item.Quantity)
		inventory.UpdateStatus()

		if err := s.repo.UpdateInventory(ctx, inventory); err != nil {
			return &inventorypb.ReleaseStockResponse{
				Success: false,
				Message: fmt.Sprintf("failed to release stock for product %s: %v", item.ProductId, err),
			}, nil
		}
	}

	return &inventorypb.ReleaseStockResponse{
		Success: true,
		Message: "stock released successfully",
	}, nil
}

// CheckStockAvailability 检查库存可用性
func (s *Server) CheckStockAvailability(ctx context.Context, req *inventorypb.CheckStockRequest) (*inventorypb.CheckStockResponse, error) {
	inventory, err := s.repo.GetInventoryByProductID(ctx, req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "inventory not found for product %s", req.ProductId)
	}

	available := inventory.Quantity - inventory.Reserved
	return &inventorypb.CheckStockResponse{
		Available:         available >= int(req.Quantity),
		AvailableQuantity: int32(available),
		ProductId:         req.ProductId,
	}, nil
}

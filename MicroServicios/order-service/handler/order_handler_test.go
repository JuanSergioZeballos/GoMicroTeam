package handler

import (
	"context"
	"errors"
	"testing"

	pb "ecommerce-microservices/proto/order"
	pb_product "ecommerce-microservices/proto/product"

	"go-micro.dev/v4/client"
)

// MockProductClient implements pb_product.ProductService for testing.
type MockProductClient struct {
	available          bool
	unavailableProduct string
	updateSuccess      bool
	checkStockCalled   bool
	updateStockCalled  bool
}

func (m *MockProductClient) GetProduct(ctx context.Context, in *pb_product.GetProductRequest, opts ...client.CallOption) (*pb_product.GetProductResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *MockProductClient) ListProducts(ctx context.Context, in *pb_product.ListProductsRequest, opts ...client.CallOption) (*pb_product.ListProductsResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *MockProductClient) CheckStock(ctx context.Context, in *pb_product.CheckStockRequest, opts ...client.CallOption) (*pb_product.CheckStockResponse, error) {
	m.checkStockCalled = true
	var unavailable []string
	if !m.available {
		unavailable = append(unavailable, m.unavailableProduct)
	}
	return &pb_product.CheckStockResponse{
		Available:           m.available,
		UnavailableProducts: unavailable,
		Message:             "mock check",
	}, nil
}

func (m *MockProductClient) UpdateStock(ctx context.Context, in *pb_product.UpdateStockRequest, opts ...client.CallOption) (*pb_product.UpdateStockResponse, error) {
	m.updateStockCalled = true
	return &pb_product.UpdateStockResponse{
		Success: m.updateSuccess,
		Message: "mock update",
		Product: &pb_product.Product{
			Id:    in.ProductId,
			Stock: 10,
		},
	}, nil
}

func TestOrderHandler(t *testing.T) {
	mockProduct := &MockProductClient{
		available:     true,
		updateSuccess: true,
	}

	h := NewOrderHandler(mockProduct)

	// Verify preloaded orders
	if len(h.orders) != 3 {
		t.Errorf("Expected 3 preloaded orders, got %d", len(h.orders))
	}

	// Test GetOrder
	var getRsp pb.GetOrderResponse
	err := h.GetOrder(context.Background(), &pb.GetOrderRequest{OrderId: "ORD-001"}, &getRsp)
	if err != nil {
		t.Fatalf("Unexpected error on GetOrder: %v", err)
	}
	if !getRsp.Found {
		t.Error("Expected order ORD-001 to be found")
	}
	if getRsp.Order.Id != "ORD-001" {
		t.Errorf("Expected order ID ORD-001, got %s", getRsp.Order.Id)
	}

	// Test ListOrders
	var listRsp pb.ListOrdersResponse
	err = h.ListOrders(context.Background(), &pb.ListOrdersRequest{UserId: "1"}, &listRsp)
	if err != nil {
		t.Fatalf("Unexpected error on ListOrders: %v", err)
	}
	if len(listRsp.Orders) != 2 {
		t.Errorf("Expected 2 orders for user 1, got %d", len(listRsp.Orders))
	}

	// Test CreateOrder Success
	var createRsp pb.CreateOrderResponse
	err = h.CreateOrder(context.Background(), &pb.CreateOrderRequest{
		UserId: "1",
		Items: []*pb.OrderItem{
			{ProductId: "P1", ProductName: "Laptop Gamer", Quantity: 1, UnitPrice: 1200.0},
		},
	}, &createRsp)

	if err != nil {
		t.Fatalf("Unexpected error on CreateOrder: %v", err)
	}
	if !createRsp.Success {
		t.Errorf("Expected CreateOrder success, message: %s", createRsp.Message)
	}
	if createRsp.Order == nil {
		t.Fatal("Expected order to be returned in CreateOrderResponse")
	}
	if createRsp.Order.Total != 1200.0 {
		t.Errorf("Expected total to be 1200.0, got %f", createRsp.Order.Total)
	}
	if !mockProduct.checkStockCalled || !mockProduct.updateStockCalled {
		t.Error("Expected CheckStock and UpdateStock to be called on ProductService client")
	}

	// Test UpdateOrderStatus
	var updateRsp pb.UpdateOrderStatusResponse
	err = h.UpdateOrderStatus(context.Background(), &pb.UpdateOrderStatusRequest{
		OrderId:   createRsp.Order.Id,
		NewStatus: "CONFIRMED",
	}, &updateRsp)
	if err != nil {
		t.Fatalf("Unexpected error on UpdateOrderStatus: %v", err)
	}
	if !updateRsp.Success {
		t.Errorf("Expected status update success, message: %s", updateRsp.Message)
	}
	if updateRsp.Order.Status != "CONFIRMED" {
		t.Errorf("Expected status CONFIRMED, got %s", updateRsp.Order.Status)
	}
}

func TestCreateOrder_OutOfStock(t *testing.T) {
	mockProduct := &MockProductClient{
		available:          false,
		unavailableProduct: "P1",
	}

	h := NewOrderHandler(mockProduct)

	var createRsp pb.CreateOrderResponse
	err := h.CreateOrder(context.Background(), &pb.CreateOrderRequest{
		UserId: "1",
		Items: []*pb.OrderItem{
			{ProductId: "P1", ProductName: "Laptop Gamer", Quantity: 1, UnitPrice: 1200.0},
		},
	}, &createRsp)

	if err != nil {
		t.Fatalf("Unexpected error on CreateOrder: %v", err)
	}
	if createRsp.Success {
		t.Error("Expected CreateOrder to fail due to stock unavailability")
	}
	if mockProduct.updateStockCalled {
		t.Error("Expected UpdateStock not to be called when product is out of stock")
	}
}

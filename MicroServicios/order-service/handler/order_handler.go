// Package handler - Order Service Handler
//
// Implementa la lógica de negocio del OrderService definida en order.proto.
// Mantiene un almacén en memoria para los pedidos (puede reemplazarse por una BD).
// Se comunica con ProductService para validar stock antes de confirmar pedidos.
//
// Implementación a cargo de: Persona 3
package handler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	pb "ecommerce-microservices/proto/order"
	pb_product "ecommerce-microservices/proto/product"
)

// OrderHandler implementa la interfaz OrderServiceHandler generada por protobuf.
type OrderHandler struct {
	mu            sync.RWMutex
	orders        map[string]*pb.Order
	nextID        int
	productClient pb_product.ProductService
}

// NewOrderHandler crea un nuevo OrderHandler con datos de prueba precargados.
// Recibe un cliente de ProductService para validar stock al crear pedidos.
func NewOrderHandler(productClient pb_product.ProductService) *OrderHandler {
	h := &OrderHandler{
		orders:        make(map[string]*pb.Order),
		nextID:        100,
		productClient: productClient,
	}

	// Datos de prueba precargados
	now := time.Now().Format(time.RFC3339)

	h.orders["ORD-001"] = &pb.Order{
		Id:     "ORD-001",
		UserId: "1",
		Items: []*pb.OrderItem{
			{ProductId: "P1", ProductName: "Laptop Gamer", Quantity: 1, UnitPrice: 1200.0},
			{ProductId: "P2", ProductName: "Mouse Wireless", Quantity: 2, UnitPrice: 40.0},
		},
		Total:     1280.0,
		Status:    "CONFIRMED",
		CreatedAt: now,
		UpdatedAt: now,
	}

	h.orders["ORD-002"] = &pb.Order{
		Id:     "ORD-002",
		UserId: "1",
		Items: []*pb.OrderItem{
			{ProductId: "P3", ProductName: "Teclado Mecánico", Quantity: 1, UnitPrice: 90.0},
		},
		Total:     90.0,
		Status:    "PENDING",
		CreatedAt: now,
		UpdatedAt: now,
	}

	h.orders["ORD-003"] = &pb.Order{
		Id:     "ORD-003",
		UserId: "2",
		Items: []*pb.OrderItem{
			{ProductId: "P2", ProductName: "Mouse Wireless", Quantity: 3, UnitPrice: 40.0},
		},
		Total:     120.0,
		Status:    "SHIPPED",
		CreatedAt: now,
		UpdatedAt: now,
	}

	log.Println("OrderHandler inicializado con 3 pedidos de prueba")
	return h
}

// CreateOrder crea un nuevo pedido.
// Valida stock con ProductService antes de confirmar.
// Flujo:
//  1. Construir CheckStockRequest a partir de los items
//  2. Llamar a ProductService.CheckStock() vía gRPC
//  3. Si stock no disponible → retornar error
//  4. Decrementar stock llamando a ProductService.UpdateStock() para cada item
//  5. Calcular total, generar ID y guardar pedido con estado PENDING
func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest, rsp *pb.CreateOrderResponse) error {
	log.Printf("CreateOrder: user_id=%s, items=%d", req.UserId, len(req.Items))

	// Validar que hay al menos un item
	if len(req.Items) == 0 {
		rsp.Success = false
		rsp.Message = "El pedido debe tener al menos un item"
		return nil
	}

	// Validar user_id
	if req.UserId == "" {
		rsp.Success = false
		rsp.Message = "El user_id es obligatorio"
		return nil
	}

	// 1. Verificar stock con ProductService
	stockItems := make([]*pb_product.StockItem, 0, len(req.Items))
	for _, item := range req.Items {
		stockItems = append(stockItems, &pb_product.StockItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	stockResp, err := h.productClient.CheckStock(ctx, &pb_product.CheckStockRequest{
		Items: stockItems,
	})
	if err != nil {
		log.Printf("Error al verificar stock: %v", err)
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Error al verificar stock: %v", err)
		return nil
	}

	if !stockResp.Available {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Stock insuficiente para productos: %s",
			strings.Join(stockResp.UnavailableProducts, ", "))
		return nil
	}

	// 2. Decrementar stock de cada producto
	for _, item := range req.Items {
		updateResp, err := h.productClient.UpdateStock(ctx, &pb_product.UpdateStockRequest{
			ProductId:      item.ProductId,
			QuantityChange: -item.Quantity, // negativo = decremento
		})
		if err != nil {
			log.Printf("Error al actualizar stock de %s: %v", item.ProductId, err)
			rsp.Success = false
			rsp.Message = fmt.Sprintf("Error al actualizar stock del producto %s: %v", item.ProductId, err)
			return nil
		}
		if !updateResp.Success {
			rsp.Success = false
			rsp.Message = fmt.Sprintf("No se pudo actualizar stock de %s: %s", item.ProductId, updateResp.Message)
			return nil
		}
	}

	// 3. Calcular total
	// NOTA: En una versión productiva, los precios deben validarse contra el ProductService
	// para evitar manipulaciones desde el cliente/frontend.
	var total float64
	for _, item := range req.Items {
		total += float64(item.Quantity) * item.UnitPrice
	}

	// 4. Crear el pedido
	h.mu.Lock()
	defer h.mu.Unlock()

	h.nextID++
	orderID := fmt.Sprintf("ORD-%03d", h.nextID)
	now := time.Now().Format(time.RFC3339)

	order := &pb.Order{
		Id:        orderID,
		UserId:    req.UserId,
		Items:     req.Items,
		Total:     total,
		Status:    "PENDING",
		CreatedAt: now,
		UpdatedAt: now,
	}
	h.orders[orderID] = order

	log.Printf("Pedido creado: %s (total: %.2f, status: PENDING)", orderID, total)

	rsp.Order = order
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Pedido %s creado exitosamente", orderID)
	return nil
}

// ListOrders devuelve todos los pedidos, opcionalmente filtrados por user_id.
// Si user_id está vacío, retorna todos los pedidos.
func (h *OrderHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest, rsp *pb.ListOrdersResponse) error {
	log.Printf("ListOrders: user_id=%s", req.UserId)

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, order := range h.orders {
		if req.UserId == "" || order.UserId == req.UserId {
			rsp.Orders = append(rsp.Orders, order)
		}
	}

	log.Printf("ListOrders: retornando %d pedidos", len(rsp.Orders))
	return nil
}

// UpdateOrderStatus actualiza el estado de un pedido existente.
// Valida que el pedido exista y que el nuevo estado sea válido.
// Estados válidos: PENDING, CONFIRMED, SHIPPED, DELIVERED, CANCELLED
func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest, rsp *pb.UpdateOrderStatusResponse) error {
	log.Printf("UpdateOrderStatus: order_id=%s, new_status=%s", req.OrderId, req.NewStatus)

	// Validar que el estado sea válido
	validStatuses := map[string]bool{
		"PENDING":   true,
		"CONFIRMED": true,
		"SHIPPED":   true,
		"DELIVERED": true,
		"CANCELLED": true,
	}

	newStatus := strings.ToUpper(req.NewStatus)
	if !validStatuses[newStatus] {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Estado inválido: %s. Estados válidos: PENDING, CONFIRMED, SHIPPED, DELIVERED, CANCELLED", req.NewStatus)
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	order, exists := h.orders[req.OrderId]
	if !exists {
		rsp.Success = false
		rsp.Message = fmt.Sprintf("Pedido %s no encontrado", req.OrderId)
		return nil
	}

	oldStatus := order.Status
	order.Status = newStatus
	order.UpdatedAt = time.Now().Format(time.RFC3339)

	log.Printf("Pedido %s actualizado: %s → %s", req.OrderId, oldStatus, newStatus)

	rsp.Order = order
	rsp.Success = true
	rsp.Message = fmt.Sprintf("Estado actualizado de %s a %s", oldStatus, newStatus)
	return nil
}

// GetOrder obtiene un pedido específico por su ID.
func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest, rsp *pb.GetOrderResponse) error {
	log.Printf("GetOrder: order_id=%s", req.OrderId)

	h.mu.RLock()
	defer h.mu.RUnlock()

	if order, ok := h.orders[req.OrderId]; ok {
		rsp.Order = order
		rsp.Found = true
		log.Printf("GetOrder: pedido %s encontrado (status: %s)", req.OrderId, order.Status)
	} else {
		rsp.Found = false
		log.Printf("GetOrder: pedido %s no encontrado", req.OrderId)
	}
	return nil
}

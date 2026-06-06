// Package handler - Order Service Handler
//
// Implementa la lógica de negocio del OrderService definida en order.proto.
// Mantiene un almacén en memoria para los pedidos (puede reemplazarse por una BD).
//
// Implementación a cargo de: Persona 3
package handler

// OrderHandler implementa la interfaz OrderServiceHandler generada por protobuf.
type OrderHandler struct {
	// TODO (Persona 3): Agregar campos necesarios
	// - orders map[string]*pb.Order  (almacén en memoria)
	// - productClient pb.ProductService (cliente para validar stock)
}

// CreateOrder crea un nuevo pedido.
// Debe llamar a ProductService.CheckStock antes de confirmar.
// func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest, rsp *pb.CreateOrderResponse) error {
//     // 1. Validar stock llamando a ProductService
//     // 2. Calcular total
//     // 3. Crear el pedido con estado PENDING
//     // 4. Guardar en el almacén
//     return nil
// }

// ListOrders devuelve todos los pedidos, opcionalmente filtrados por user_id.
// func (h *OrderHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest, rsp *pb.ListOrdersResponse) error {
//     return nil
// }

// UpdateOrderStatus actualiza el estado de un pedido existente.
// func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest, rsp *pb.UpdateOrderStatusResponse) error {
//     return nil
// }

// GetOrder obtiene un pedido específico por su ID.
// func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest, rsp *pb.GetOrderResponse) error {
//     return nil
// }

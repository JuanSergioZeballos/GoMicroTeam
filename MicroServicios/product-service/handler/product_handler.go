// Package handler - Product Service Handler
//
// Implementa la lógica de negocio del ProductService definida en product.proto.
// Mantiene un catálogo de productos en memoria con datos de prueba precargados.
//
// Implementación a cargo de: Persona 4
package handler

import (
	"context"
	"sync"

	pb "ecommerce-microservices/proto/product"
)
// ProductHandler implementa la interfaz ProductServiceHandler generada por protobuf.
type ProductHandler struct {
	// TODO (Persona 4): Agregar campos necesarios
	// - products map[string]*pb.Product  (catálogo en memoria)
	mu       sync.RWMutex
	products map[string]*pb.Product
}

func NewProductHandler() *ProductHandler {
	h := &ProductHandler{
		products: make(map[string]*pb.Product),
	}

	h.products["P1"] = &pb.Product{Id: "P1", Name: "Laptop Gamer", Description: "i7, 16GB RAM, RTX 3060", Price: 1200.0, Stock: 10, Category: "Electronics"}
	h.products["P2"] = &pb.Product{Id: "P2", Name: "Mouse Wireless", Description: "Logitech G305", Price: 40.0, Stock: 50, Category: "Peripherals"}
	h.products["P3"] = &pb.Product{Id: "P3", Name: "Teclado Mecánico", Description: "Keychron K2", Price: 90.0, Stock: 20, Category: "Peripherals"}

	return h
}

// GetProduct obtiene un producto por su ID.
// func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest, rsp *pb.GetProductResponse) error {
//     return nil
// }

func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest, rsp *pb.GetProductResponse) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if product, ok := h.products[req.ProductId]; ok {
		rsp.Product = product
		rsp.Found = true
	} else {
		rsp.Found = false
	}
	return nil
}

// ListProducts lista todos los productos, opcionalmente filtrados por categoría.
// func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest, rsp *pb.ListProductsResponse) error {
//     return nil
// }

func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest, rsp *pb.ListProductsResponse) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, p := range h.products {
		if req.Category == "" || p.Category == req.Category {
			rsp.Products = append(rsp.Products, p)
		}
	}
	return nil
}

// CheckStock verifica si hay stock suficiente para una lista de ítems.
// func (h *ProductHandler) CheckStock(ctx context.Context, req *pb.CheckStockRequest, rsp *pb.CheckStockResponse) error {
//     return nil
// }

func (h *ProductHandler) CheckStock(ctx context.Context, req *pb.CheckStockRequest, rsp *pb.CheckStockResponse) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	rsp.Available = true
	for _, item := range req.Items {
		p, ok := h.products[item.ProductId]
		if !ok || p.Stock < item.Quantity {
			rsp.Available = false
			rsp.UnavailableProducts = append(rsp.UnavailableProducts, item.ProductId)
		}
	}

	if !rsp.Available {
		rsp.Message = "Stock insuficiente para algunos productos"
	} else {
		rsp.Message = "Stock disponible"
	}
	return nil
}

// UpdateStock actualiza el stock de un producto (incremento o decremento).
// func (h *ProductHandler) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest, rsp *pb.UpdateStockResponse) error {
//     return nil
// }

func (h *ProductHandler) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest, rsp *pb.UpdateStockResponse) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	p, ok := h.products[req.ProductId]
	if !ok {
		rsp.Success = false
		rsp.Message = "Producto no encontrado"
		return nil
	}

	newStock := p.Stock + req.QuantityChange
	if newStock < 0 {
		rsp.Success = false
		rsp.Message = "No se puede reducir el stock por debajo de 0"
		return nil
	}

	p.Stock = newStock
	rsp.Product = p
	rsp.Success = true
	rsp.Message = "Stock actualizado correctamente"
	return nil
}


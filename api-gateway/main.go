// Package main - API Gateway
//
// Punto de entrada HTTP/REST que expone los microservicios internos al exterior.
// Utiliza Gin como router HTTP y Go Micro como cliente para comunicarse
// con los servicios internos (Order, Product, User) vía gRPC.
//
// Responsabilidades:
//   - Exponer endpoints REST para clientes externos (frontend, mobile)
//   - Enrutar peticiones al servicio interno correspondiente
//   - Middleware de autenticación (validar token del UserService)
//   - Agregación de respuestas cuando sea necesario
//
// Endpoints planificados:
//   POST   /api/v1/auth/register      → UserService.Register
//   POST   /api/v1/auth/login          → UserService.Login
//   GET    /api/v1/users/:id           → UserService.GetUser
//   GET    /api/v1/products            → ProductService.ListProducts
//   GET    /api/v1/products/:id        → ProductService.GetProduct
//   POST   /api/v1/orders              → OrderService.CreateOrder
//   GET    /api/v1/orders              → OrderService.ListOrders
//   GET    /api/v1/orders/:id          → OrderService.GetOrder
//   PATCH  /api/v1/orders/:id/status   → OrderService.UpdateOrderStatus
//
// Implementación a cargo de: Persona 5
package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("=== API Gateway ===")
	fmt.Println("Gateway HTTP/REST → microservicios gRPC")
	fmt.Println("Puerto HTTP: 8080")

	// TODO (Persona 5): Implementar el gateway
	//
	// Pasos esperados:
	// 1. Crear servicio Go Micro (web.NewService o micro.NewService)
	// 2. Crear clientes Go Micro para cada servicio:
	//    - orderClient  = pb_order.NewOrderService("order-service", service.Client())
	//    - productClient = pb_product.NewProductService("product-service", service.Client())
	//    - userClient   = pb_user.NewUserService("user-service", service.Client())
	// 3. Configurar router Gin con los endpoints listados arriba
	// 4. Iniciar en puerto 8080

	log.Println("API Gateway - stub listo para implementación por Persona 5")
}

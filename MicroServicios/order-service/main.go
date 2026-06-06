// Package main - Order Service
//
// Microservicio de gestión de pedidos para la plataforma de comercio electrónico.
// Utiliza Go Micro v4 como framework de microservicios.
//
// Responsabilidades:
//   - Crear nuevos pedidos (validando stock con ProductService)
//   - Listar pedidos por usuario o globalmente
//   - Actualizar el estado de un pedido (PENDING → CONFIRMED → SHIPPED → DELIVERED)
//   - Consultar un pedido específico por ID
//
// Este servicio se registra automáticamente en el service registry (etcd)
// y expone sus endpoints mediante gRPC a través de Go Micro.
//
// Implementación a cargo de: Persona 3
package main

import (
	"fmt"
	"log"

	// Descomentar cuando se instalen las dependencias:
	// "go-micro.dev/v4"
	// pb "ecommerce-microservices/proto/order"
)

func main() {
	fmt.Println("=== Order Service ===")
	fmt.Println("Servicio de gestión de pedidos")
	fmt.Println("Puerto gRPC: se asigna dinámicamente vía Go Micro")

	// TODO (Persona 3): Implementar el servidor Go Micro
	//
	// Pasos esperados:
	// 1. Crear instancia de micro.NewService con nombre "order-service"
	// 2. Registrar el handler OrderServiceHandler
	// 3. Iniciar el servicio con service.Run()
	//
	// Ejemplo de estructura:
	//
	//   service := micro.NewService(
	//       micro.Name("order-service"),
	//       micro.Version("1.0.0"),
	//   )
	//   service.Init()
	//   pb.RegisterOrderServiceHandler(service.Server(), &handler.OrderHandler{})
	//   if err := service.Run(); err != nil {
	//       log.Fatal(err)
	//   }

	log.Println("Order Service - stub listo para implementación por Persona 3")
}

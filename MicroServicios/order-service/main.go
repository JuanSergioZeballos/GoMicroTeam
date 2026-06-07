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

	"go-micro.dev/v4"
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	_ "github.com/go-micro/plugins/v4/client/grpc"
	_ "github.com/go-micro/plugins/v4/server/grpc"
	_ "github.com/go-micro/plugins/v4/transport/grpc"

	"ecommerce-microservices/MicroServicios/order-service/handler"
	pb "ecommerce-microservices/proto/order"
	pb_product "ecommerce-microservices/proto/product"
)

func main() {
	fmt.Println("=== Order Service ===")
	fmt.Println("Servicio de gestión de pedidos")

	// 1. Crear instancia de micro.NewService con nombre "order-service"
	service := micro.NewService()
	service.Init(
		micro.Name("order-service"),
		micro.Version("1.0.0"),
	)

	// 3. Crear cliente de ProductService para validar stock
	productClient := pb_product.NewProductService("product-service", service.Client())

	// 4. Registrar el handler OrderServiceHandler con el cliente de ProductService inyectado
	pb.RegisterOrderServiceHandler(service.Server(), handler.NewOrderHandler(productClient))

	log.Println("Order Service iniciado - registrado en service registry")

	// 5. Iniciar el servicio
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

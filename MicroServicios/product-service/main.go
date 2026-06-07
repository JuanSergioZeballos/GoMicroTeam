// Package main - Product Service
//
// Microservicio de gestión de catálogo e inventario de productos.
// Utiliza Go Micro v4 como framework de microservicios.
//
// Responsabilidades:
//   - Consultar productos del catálogo
//   - Verificar disponibilidad de stock
//   - Actualizar stock al confirmar pedidos
//
// Este servicio se registra en el service registry (etcd)
// y es consumido principalmente por OrderService para validar stock.
//
// Implementación a cargo de: Persona 4
package main

import (
	"log"
    "fmt"
	"go-micro.dev/v4"
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	_ "github.com/go-micro/plugins/v4/client/grpc"
	_ "github.com/go-micro/plugins/v4/server/grpc"
	_ "github.com/go-micro/plugins/v4/transport/grpc"
	"ecommerce-microservices/MicroServicios/product-service/handler"
	pb "ecommerce-microservices/proto/product"
)

func main() {
	fmt.Println("=== Product Service ===")
	fmt.Println("Servicio de catálogo e inventario")
	service := micro.NewService()
	service.Init(
		micro.Name("product-service"),
		micro.Version("1.0.0"),
	)

	pb.RegisterProductServiceHandler(service.Server(), handler.NewProductHandler())

	log.Println("Product Service iniciado - registrado en service registry")

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

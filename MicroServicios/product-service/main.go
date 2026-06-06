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
	"ecommerce-microservices/MicroServicios/product-service/handler"
	pb "ecommerce-microservices/proto/product"
)

func main() {
	fmt.Println("=== Product Service ===")
	fmt.Println("Servicio de catálogo e inventario")
	service := micro.NewService(
		micro.Name("product-service"),
		micro.Version("1.0.0"),
	)

	service.Init()

	pb.RegisterProductServiceHandler(service.Server(), handler.NewProductHandler())

	log.Println("Product Service iniciado - registrado en service registry")

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

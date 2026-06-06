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
	"fmt"
	"log"
)

func main() {
	fmt.Println("=== Product Service ===")
	fmt.Println("Servicio de catálogo e inventario")

	// TODO (Persona 4): Implementar el servidor Go Micro
	//
	//   service := micro.NewService(
	//       micro.Name("product-service"),
	//       micro.Version("1.0.0"),
	//   )
	//   service.Init()
	//   pb.RegisterProductServiceHandler(service.Server(), &handler.ProductHandler{})
	//   if err := service.Run(); err != nil {
	//       log.Fatal(err)
	//   }

	log.Println("Product Service - stub listo para implementación por Persona 4")
}

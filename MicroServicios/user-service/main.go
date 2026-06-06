// Package main - User Service
//
// Microservicio de gestión de usuarios y autenticación básica.
// Utiliza Go Micro v4 como framework de microservicios.
//
// Responsabilidades:
//   - Registro de nuevos usuarios
//   - Autenticación (login) con generación de token
//   - Consulta de perfil de usuario
//   - Listado de usuarios registrados
//
// Este servicio se registra en el service registry (etcd)
// y es consumido por el API Gateway para autenticación.
//
// Implementación a cargo de: Persona 4
package main

import (
	"fmt"
	"log"
	"go-micro.dev/v4"
	"ecommerce-microservices/MicroServicios/user-service/handler"
	pb "ecommerce-microservices/proto/user"
)

func main() {
	fmt.Println("=== User Service ===")
	fmt.Println("Servicio de usuarios y autenticación")

	// TODO (Persona 4): Implementar el servidor Go Micro
	//
	//   service := micro.NewService(
	//       micro.Name("user-service"),
	//       micro.Version("1.0.0"),
	//   )
	//   service.Init()
	//   pb.RegisterUserServiceHandler(service.Server(), &handler.UserHandler{})
	//   if err := service.Run(); err != nil {
	//       log.Fatal(err)
	//   }
	service := micro.NewService(
		micro.Name("user-service"),
		micro.Version("1.0.0"),
	)

	service.Init()

	pb.RegisterUserServiceHandler(service.Server(), handler.NewUserHandler())

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("User Service - stub listo para implementación por Persona 4")
}
